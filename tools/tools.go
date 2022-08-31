package tools

import (
	"database/sql"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var IsProduction bool
var ServerID int32
var MaxUploadFileSize int64
var MaxUploadedFilesSize int64
var MaxUploadedFiles int64
var MaxSessionsFromClient = 7 //if config value is empty
var ServerAddr string
var ServerPort string
var FileStorage = "/storage"
var ImageStorage = FileStorage + "/images"
var MaxAvatarSize = 128
var MinAvatarSize = 40
var MaxVertexWidthSize = 450
var MinVertexWidthSize = 250

const EmptySqlResultSet = "sql: no rows in result set"

type ImageMode int

const (
	AvatarMode ImageMode = 1
	VertexMode ImageMode = 2
)

func CheckLang(lang string) bool {

	if lang == "en" || lang == "ru" {
		return true
	}

	return false
}

//You can use it instead of db or tx if you don't know type of param
type QueryExecutor interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var mailcheck *regexp.Regexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var fileExtensionCheck = regexp.MustCompile("(png|svg|jpg|jpeg)$")

func ValidateEmail(email string) bool {
	return mailcheck.MatchString(email)
}

func ValidateImageExtension(fileName string) bool {
	return fileExtensionCheck.MatchString(fileName)
}

func SqlUint64Seq(nums []uint64) string { // [1, 2, 3, 4] => 1,2,3,4
	if len(nums) == 0 {
		return "0"
	}

	b := make([]byte, 0, len(nums)*9)
	for _, n := range nums {
		b = strconv.AppendInt(b, int64(n), 10)
		b = append(b, ',')
	}
	b = b[:len(b)-1]
	return string(b)
}

func SqlUint32Seq(nums []uint32) string { // [1, 2, 3, 4] => 1,2,3,4
	if len(nums) == 0 {
		return "0"
	}

	b := make([]byte, 0, len(nums)*9)
	for _, n := range nums {
		b = strconv.AppendInt(b, int64(n), 10)
		b = append(b, ',')
	}
	b = b[:len(b)-1]
	return string(b)
}

func SqlUint64BracketsSeq(nums []uint64) string { // [1, 2, 3, 4] => (1),(2),(3),(4)

	if len(nums) == 0 {
		return "(0)"
	}

	estimate := len(nums) * 4
	b := make([]byte, 0, estimate)
	b = append(b, '(')
	for _, n := range nums {
		b = strconv.AppendInt(b, int64(n), 10)
		b = append(b, "),("...)
	}
	b = b[:len(b)-2]
	return string(b)
}

func IsUint64NumInSlice(num uint64, slice []uint64) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == num {
			return true
		}
	}
	return false
}

func IsAdmin(actorId uint32) bool {
	count := 0
	err := DBS.QueryRow("SELECT count(id) FROM actor WHERE id = $1 AND admin > 0", actorId).Scan(&count)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if count != 0 {
		return true
	} else {
		return false
	}
}

func IsUint32NumInSlice(num uint32, slice []uint32) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == num {
			return true
		}
	}
	return false
}

func IsInt64NumInSlice(num int64, slice []int64) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == num {
			return true
		}
	}
	return false
}

func IsDueDateTimeValid(needleTime time.Time) bool {
	currentTime := time.Now()

	if needleTime.Unix() != -62135596800 && needleTime.Year() < currentTime.Year() {
		return false
	}

	if needleTime.Unix() != -62135596800 && needleTime.Year() == currentTime.Year() && needleTime.Month() < currentTime.Month() {
		return false
	}

	if needleTime.Unix() != -62135596800 && needleTime.Month() == currentTime.Month() && needleTime.Day() < currentTime.Day() {
		return false
	}

	return true
}

func GetTimeForDB(time time.Time) *string {
	//if time is equal to 0001-01-01 00:00:00 then nil
	if time.Unix() == -62135596800 {
		return nil
	}
	returnTime := time.Format("2006-01-02 15:04:05")
	return &returnTime
}

func GetUnixTimeToDB(timeToUnix int64) sql.NullTime {
	if timeToUnix == 0 {
		return sql.NullTime{Valid: false}
	} else {
		return sql.NullTime{Valid: true, Time: time.Unix(timeToUnix, 0)}
	}
}

var DBS *sql.DB

func ModifyHTTPCors(next http.Handler) http.Handler {

	var status string
	var url string

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if IsProduction {
			res.Header().Set("Access-Control-Allow-Origin", ServerAddr)
		} else {
			res.Header().Set("Access-Control-Allow-Origin", ServerAddr+":4321")
		}
		url = req.URL.Path
		res.Header().Set("Access-Control-Allow-Headers", "X-Session-ID, Content-Type, *")
		res.Header().Set("Access-Control-Allow-Credentials", "true")
		res.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTION, PATCH, DELETE")
		status = "success"
		HttpRequestsCounter.WithLabelValues(status, url).Inc()
		next.ServeHTTP(res, req)
	})
}

func Time2DB(t time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

var reg = regexp.MustCompile(`(\n|\r|\a|\f|\t|\v)`)

func CheckVariableOrCreate(token string, label string) (uint64, uint64) {
	tx, err := DBS.Begin()
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
	}

	var userId uint64
	err = tx.QueryRow("SELECT id FROM users WHERE token = $1", token).Scan(&userId)
	if err != nil {
		fmt.Println("aa", err)
		tx.Rollback()
		return 0, 0
	}

	var varId uint64
	err = tx.QueryRow("SELECT uv.varid FROM users_variables uv INNER JOIN variables v ON v.id = uv.varid WHERE v.label = $1", label).Scan(&varId)
	if err != nil && err.Error() != EmptySqlResultSet {
		tx.Rollback()
		return 0, 0
	}

	if varId == 0 {
		tx.QueryRow("INSERT INTO variables(name, label) VALUES($1, $2) RETURNING id", label, label).Scan(&varId)
		tx.Exec("INSERT INTO users_variables(userid, varid) VALUES($1, $2)", userId, varId)
	}

	tx.Commit()

	return userId, varId
}
func commaShielding(str string) string {
	return strings.Replace(str, "'", "''", 1)
}

func GetIntParamFromRequestUrl(r *http.Request, paramKey string) (param int, ok bool) {
	if len(r.URL.Query()[paramKey]) > 0 && len(r.URL.Query()[paramKey][0]) != 0 {
		param, _ := strconv.Atoi(r.URL.Query()[paramKey][0])
		return param, true
	}
	return 0, false
}

func GetStringParamFromRequestUrl(r *http.Request, paramKey string) (param string, ok bool) {
	if len(r.URL.Query()[paramKey]) > 0 && len(r.URL.Query()[paramKey][0]) != 0 {
		return r.URL.Query()[paramKey][0], true
	}
	return "", false
}

var HttpRequestsCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_request_counter", // metric name
		Help: "Number of get_books request.",
	},
	[]string{"status", "url"}, // labels
)

func init() {
	prometheus.MustRegister(HttpRequestsCounter)
}
