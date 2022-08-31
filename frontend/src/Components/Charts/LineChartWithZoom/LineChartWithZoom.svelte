<style>
        #chartDiv {
                width: 100%;
                height: 500px;
        }
</style>
<script>
import * as am5 from "@amcharts/amcharts5";
import * as am5xy from "@amcharts/amcharts5/xy";
import am5themes_Animated from "@amcharts/amcharts5/themes/Animated";
import ru from "@amcharts/amcharts5/locales/ru_RU"
import { onMount } from 'svelte';
/* Chart code */
// Create root element
// https://www.amcharts.com/docs/v5/getting-started/#Root_element

export let data;

let chartDiv;

onMount(()=>{
        data.map((val)=>{
            val.Current_Value = parseInt(val.Current_Value)
        })
        console.log(data);
        let root = am5.Root.new(chartDiv);

// Set themes
// https://www.amcharts.com/docs/v5/concepts/themes/
        root.setThemes([
                am5themes_Animated.new(root)
        ]);
        root.locale = ru
// Create chart
// https://www.amcharts.com/docs/v5/charts/xy-chart/
        let chart = root.container.children.push(am5xy.XYChart.new(root, {
                panX: true,
                panY: true,
                wheelX: "panX",
                wheelY: "zoomX",
                pinchZoomX:true,
        }));




// Create axes
// https://www.amcharts.com/docs/v5/charts/xy-chart/axes/
        let xAxis = chart.xAxes.push(am5xy.DateAxis.new(root, {
                baseInterval: { timeUnit: "second", count: 1 },
                renderer: am5xy.AxisRendererX.new(root, {}),
                tooltip: am5.Tooltip.new(root, {})
        }));

        let yAxis = chart.yAxes.push(am5xy.ValueAxis.new(root, {
                renderer: am5xy.AxisRendererY.new(root, {})
        }));

    // Add cursor
// https://www.amcharts.com/docs/v5/charts/xy-chart/cursor/


        let tooltip = am5.Tooltip.new(root, {
            getFillFromSprite: false,
            getStrokeFromSprite: false,
            autoTextColor: false,
            getLabelFillFromSprite: false,
            labelText: "{valueY}",
        });

        tooltip.get("background").setAll({
            fill: "#1c1e33",
            fillOpacity: 0.8
        });
        tooltip.label.setAll({
            fill: "#fff"
        });

        let series = chart.series.push(am5xy.LineSeries.new(root, {
                name: "Series",
                xAxis: xAxis,
                yAxis: yAxis,
                valueYField: "Current_Value",
                valueXField: "CreationDate",
                tooltip: tooltip
        }));

        let cursor = chart.set("cursor", am5xy.XYCursor.new(root, {
            xAxis: xAxis,
            behavior: "zoomXY",
            snapToSeries: [ series ]
        }));
        cursor.lineY.set("visible", false);

        let scrollbar = chart.set("scrollbarX", am5xy.XYChartScrollbar.new(root, {
                orientation: "horizontal",
                height: 60
        }));

        let sbDateAxis = scrollbar.chart.xAxes.push(am5xy.DateAxis.new(root, {
                baseInterval: {
                        timeUnit: "second",
                        count: 1
                },
                renderer: am5xy.AxisRendererX.new(root, {})
        }));

        let sbValueAxis = scrollbar.chart.yAxes.push(
                am5xy.ValueAxis.new(root, {
                        renderer: am5xy.AxisRendererY.new(root, {})
                })
        );

        let sbSeries = scrollbar.chart.series.push(am5xy.LineSeries.new(root, {
                valueYField: "Current_Value",
                valueXField: "CreationDate",
                xAxis: sbDateAxis,
                yAxis: sbValueAxis
        }));

        series.strokes.template.setAll({
            strokeWidth: 2,
        });
        series.data.setAll(data);
        sbSeries.data.setAll(data);


        series.appear(1000);
        chart.appear(1000, 100);
})

</script>
<div id="chartDiv" bind:this={chartDiv}></div>
<!--<script>-->
<!--    import * as am4core from "@amcharts/amcharts4/core";-->
<!--    import * as am4charts from "@amcharts/amcharts4/charts";-->
<!--    import am4themes_animated from "@amcharts/amcharts4/themes/animated";-->
<!--    import am4lang_en_US from "@amcharts/amcharts4/lang/ru_RU";-->
<!--    import {afterUpdate, beforeUpdate, onMount} from "svelte";-->

<!--    export let data = [];-->
<!--    export let style = "";-->


<!--    let paddingRight;-->
<!--    let prevDataLength;-->
<!--    afterUpdate(()=>{-->
<!--        console.log("update")-->
<!--            if (prevDataLength && prevDataLength !== data.length) {-->
<!--                console.log("push")-->
<!--                console.log(prevDataLength - 1, data.length);-->
<!--                let i = prevDataLength - 1;-->
<!--                let end = data.length;-->
<!--                while (i < end) {-->
<!--                    chart.addData(data[i]);-->
<!--                    i++-->
<!--                }-->
<!--                prevDataLength = data.length;-->
<!--            }-->
<!--    })-->

<!--    let chart;-->
<!--    onMount(()=>{-->
<!--        console.log("mount");-->
<!--        am4core.useTheme(am4themes_animated);-->

<!--        chart = am4core.create("chartDiv", am4charts.XYChart);-->
<!--        chart.language.locale = am4lang_en_US;-->
<!--// Add data-->

<!--        console.log(data);-->
<!--        chart.data = data;-->
<!--        // chart.data = [{-->
<!--        //     CreationDate: "2019-05-11T21:06:08.852Z",-->
<!--        //     Current_Value:1271-->
<!--        // }, {-->
<!--        //     CreationDate: "2020-05-11T21:06:08.852Z",-->
<!--        //     Current_Value:500-->
<!--        // }];-->

<!--// Create axes-->
<!--        let dateAxis = chart.xAxes.push(new am4charts.DateAxis());-->
<!--        dateAxis.renderer.minGridDistance = 50;-->
<!--        dateAxis.tooltip.background.fill = am4core.color("#212429");-->
<!--        dateAxis.tooltip.background.strokeWidth = 0;-->
<!--        let valueAxis = chart.yAxes.push(new am4charts.ValueAxis());-->
<!--        valueAxis.tooltip.background.fill = am4core.color("#212429");-->
<!--        valueAxis.tooltip.background.strokeWidth = 0;-->
<!--// Create series-->
<!--        let series = chart.series.push(new am4charts.LineSeries());-->
<!--        series.dataFields.valueY = "Current_Value";-->
<!--        series.dataFields.dateX = "CreationDate";-->
<!--        series.strokeWidth = 2;-->
<!--        series.tensionX = 1;-->
<!--        series.minBulletDistance = 20;-->
<!--        series.tooltipText = "{valueY}";-->
<!--        series.tooltip.pointerOrientation = "vertical";-->
<!--        series.tooltip.getFillFromObject = false;-->
<!--        series.tooltip.background.fill = am4core.color("#4699d6");-->
<!--        series.tooltip.background.strokeWidth = 0;-->
<!--        series.tooltip.background.filters.clear();-->
<!--        series.tooltip.background.cornerRadius = 10;-->
<!--        series.tooltip.background.fillOpacity = 0.8;-->

<!--        series.tooltip.label.padding(12,12,12,12)-->

<!--// Add scrollbar-->
<!--        chart.scrollbarX = new am4charts.XYChartScrollbar();-->
<!--        chart.scrollbarX.series.push(series);-->

<!--// Add cursor-->
<!--        chart.cursor = new am4charts.XYCursor();-->
<!--        chart.cursor.xAxis = dateAxis;-->
<!--        chart.cursor.snapToSeries = series;-->

<!--        paddingRight = chart.paddingRight;-->

<!--        prevDataLength = data.length;-->
<!--    })-->


<!--</script>-->
<!--{#if data.length > 0}-->
<!--<div class="chartDiv" style={style}></div>-->
<!--    {/if}-->