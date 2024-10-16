function getWeeklyEventTime(arr) {
    return arr?.map(v => new Date(v.Timestamp).toDateString()) ?? [];
}

function getTotal(arr) {
    return arr?.map(v => v.Total) ?? [];
}

function getEventType(arr) {
    return arr?.map(v => v.EventType) ?? [];
}

function getEventLabel(arr) {
    const result = [];
    arr?.map(v => result.push({
        x: v.EventLabel,
        y: v.Total
    })) ?? [];
    return result;
}

const UPDATE_INTERVAL = 10000; // 10 seconds
const projectID = document.getElementById("project-id")?.textContent;
const id = projectID ? JSON.parse(projectID) : null;
const weeklyEventElem = document.getElementById("weekly-events")?.textContent;
const weeklyEvent = weeklyEventElem ? JSON.parse(weeklyEventElem) : [];
const eventTypeElem = document.getElementById("event-type-chart")?.textContent;
const eventTypes = eventTypeElem ? JSON.parse(eventTypeElem) : [];
const eventLabelElem = document.getElementById("event-label-chart")?.textContent;
const eventLabels = eventLabelElem ? JSON.parse(eventLabelElem) : [];
let weeklyEventTimestamp = getWeeklyEventTime(weeklyEvent);
let weeklyEventData = getTotal(weeklyEvent);
let eventTypesLabels = getEventType(eventTypes);
let eventTypesPercentage = getTotal(eventTypes);
let eventLabelsData = getEventLabel(eventLabels);
const weeklyEventTotalElem = document.getElementById("weekly-event-total");

const areaOptions = {
    chart: {
        height: "300px",
        maxWidth: "100%",
        type: "area",
        dropShadow: {
            enabled: true,
        },
        toolbar: {
            show: false,
        },
        animations: {
            enabled: false,
        },
        zoom: {
            enabled: false,
        }
    },
    tooltip: {
        enabled: true,
        x: {
            show: true,
        },
    },
    fill: {
        type: "gradient",
        gradient: {
            opacityFrom: 0.55,
            opacityTo: 0,
            shade: "#1C64F2",
            gradientToColors: ["#1C64F2"],
        },
    },
    dataLabels: {
        enabled: false,
        style: {
            fontFamily: "Space Grotesk",
        },
    },
    stroke: {
        width: 6,
    },
    grid: {
        show: false,
        strokeDashArray: 4,
        padding: {
            left: 2,
            right: 2,
            top: 0
        },
    },
    series: [
        {
            name: "Events",
            data: weeklyEventData,
            color: "#1A56DB",
        },
    ],
    xaxis: {
        categories: weeklyEventTimestamp,
        labels: {
            show: false,
        },
        axisBorder: {
            show: false,
        },
        axisTicks: {
            show: false,
        },
    },
    yaxis: {
        show: false,
    },
}

const pieOptions = {
    series: eventTypesPercentage,
    chart: {
        height: "380px",
        width: "100%",
        type: "pie",
        animations: {
            enabled: false,
        }
    },
    stroke: {
        colors: ["white"],
        lineCap: "",
    },
    plotOptions: {
        pie: {
            labels: {
                show: true,
            },
            size: "100%",
            dataLabels: {
                offset: -25
            }
        },
    },
    labels: eventTypesLabels,
    dataLabels: {
        enabled: true,
        style: {
            fontFamily: "Space Grotesk",
        },
    },
    legend: {
        position: "bottom",
        fontFamily: "Space Grotesk, sans-serif",
    },
    xaxis: {
        axisTicks: {
            show: true,
        },
        axisBorder: {
            show: true,
        },
    },
}

const columnOptions = {
    series: [
        {
            name: "X Event Fired",
            data: eventLabelsData,
        },
    ],
    chart: {
        type: "bar",
        height: "350px",
        fontFamily: "Space Grotesk, sans-serif",
        toolbar: {
            show: false,
        },
        animations: {
            enabled: false,
        },
    },
    plotOptions: {
        bar: {
            distributed: true,
            horizontal: false,
            columnWidth: "70%",
            borderRadiusApplication: "end",
            borderRadius: 8,
        },
    },
    tooltip: {
        shared: true,
        intersect: false,
        style: {
            fontFamily: "Space Grotesk, sans-serif",
        },
    },
    states: {
        hover: {
            filter: {
                type: "darken",
                value: 1,
            },
        },
    },
    stroke: {
        show: true,
        width: 0,
        colors: ["transparent"],
    },
    grid: {
        show: false,
        strokeDashArray: 4,
        padding: {
            left: 2,
            right: 2,
            top: -14
        },
    },
    dataLabels: {
        enabled: true,
    },
    legend: {
        show: false,
    },
    xaxis: {
        floating: true,
        labels: {
            show: false,
            style: {
                fontFamily: "Space Grotesk, sans-serif",
                cssClass: 'text-xs font-light fill-gray-500 dark:fill-gray-400'
            }
        },
        axisBorder: {
            show: true,
        },
        axisTicks: {
            show: true,
        },
    },
    yaxis: {
        show: false,
    },
    fill: {
        opacity: 1,
    },
}

const areaChart = new ApexCharts(document.getElementById("area-chart"), areaOptions);
const pieChart = new ApexCharts(document.getElementById("pie-chart"), pieOptions);
const columnChart = new ApexCharts(document.getElementById("column-chart"), columnOptions);

areaChart.render();
pieChart.render();
columnChart.render();

setInterval(function() {
    Promise.allSettled([
        updateAreaChart(areaChart),
        updatePieChart(pieChart),
        updateBarChart(columnChart),
    ])
}, UPDATE_INTERVAL);

function updateAreaChart(chart) {
    $.ajax({
        url: `/api/json/event/chart/${id}`,
        type: "GET",
        success: function(data) {
            const newSeries = getTotal(data.Time);

            // Update total text if available
            if (weeklyEventTotalElem) {
                weeklyEventTotalElem.innerText = data.Total ?? 0;
            }

            // update chart with a new data series
            chart.updateSeries([
                {
                    name: "Events",
                    data: newSeries,
                    color: "#1A56DB",
                },
            ]);
        },
        error: function() {
            console.log("Error fetching new chart data");
        }
    });
}

function updatePieChart(chart) {
    $.ajax({
        url: `/api/json/event-type/chart/${id}`,
        type: "GET",
        success: function(data) {
            const newSeries = getTotal(data);
            const newLabels = getEventType(data);

            pieOptions.labels = newLabels;
            pieOptions.series = newSeries;

            // update chart with a new data series
            chart.updateOptions(pieOptions);
        },
        error: function() {
            console.log("Error fetching new chart data");
        }
    });
}

function updateBarChart(chart) {
    $.ajax({
        url: `/api/json/event-label/chart/${id}`,
        type: "GET",
        success: function(data) {
            const newSeries = getEventLabel(data);

            columnOptions.series = [
                {
                    name: "X Event Fired",
                    data: newSeries,
                },
            ];

            // update chart with a new data series
            chart.updateOptions(columnOptions);
        },
        error: function() {
            console.log("Error fetching new chart data");
        }
    });
}
