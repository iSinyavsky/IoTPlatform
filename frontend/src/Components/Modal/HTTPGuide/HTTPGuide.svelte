<script>
    import Modal from "../index.svelte"
    import {profile} from "../../../store";
    import {variables} from "../../../store"

    let visible = true;
</script>

{#if visible}
    <Modal width="800">
        <h1>Гайд по работе с HTTP</h1>
        <div style="text-align: left; margin: 10px 30px; overflow: auto; height: 400px">
            В замену MQTT вы можете использовать стандартный протокол HTTP.
            Для доступа к устройству по HTTP, необходимо собрать уникальный URL с токеном и api label.
            <div>Пример URL: <code>siiv.ru/http-api/ваш_токен/label_устройства</code></div>
            <br>
            <div><strong>Для отправки показания</strong> вызовите POST метод, указав тело запроса в формате JSON</div>
            <code>{@html "{'value': 'значение'}"}</code>
            <br><br>
            <div><strong>Запрос показаний</strong> происходит в виде GET запроса, данные возвращаются в JSON формате</div>
            <code>{@html "{\"errno\": 0, \"count\": 0, \"data\": [{\"value\":\"40.000000\",\"time\":\"2022-03-28T13:29:38.642411Z\"}]}"}</code>
           <br><br>
            <div style="font-weight: bold">Так же в запросах можно использовать query параметры для фильтрации данных:</div>
            <br>
            <div>{@html "<strong>Выбор диапазанов по периоду: </strong> - lt (<), le(<=), eq(==), ne(!=), ge(>=), gt(>). Например lt=2022-01-01, выдаст все показания до 1 января 2022 года"}</div>
            <br>
            <div>{@html "<strong>Сортировка: </strong> order=desc/asc, по убыванию / по возрастанию"}</div>
            <br>
            <div>{@html "<strong>Ограничение данных по количеству: </strong> limit = 1...n"}</div>
            <br>
            <div>Пример запроса со всеми фильтрами</div>
            <code>siiv.ru/http-api/ваш_токен/label_устройства?lt=2022-01-01&limit=10&order=desc</code>
            <div></div>
        </div>
    </Modal>
{/if}