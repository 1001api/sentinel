function getWeeklyEventTime(arr) {
    return arr?.map(v => new Date(v.Timestamp).toDateString()) ?? [];
}

function getTotal(arr) {
    return arr?.map(v => v.Total) ?? [];
}

function getEventType(arr) {
    return arr?.map(v => v.EventType) ?? [];
}

const UPDATE_INTERVAL = 10000; // 10 seconds
const projectID = document.getElementById("project-id")?.textContent;
const id = projectID ? JSON.parse(projectID) : null;
const weeklyEventElem = document.getElementById("weekly-events")?.textContent;
const weeklyEvent = weeklyEventElem ? JSON.parse(weeklyEventElem) : [];
const eventTypeElem = document.getElementById("event-type-chart")?.textContent;
const eventTypes = eventTypeElem ? JSON.parse(eventTypeElem) : [];
let weeklyEventTimestamp = getWeeklyEventTime(weeklyEvent);
let weeklyEventData = getTotal(weeklyEvent);
let eventTypesLabels = getEventType(eventTypes);
let eventTypesPercentage = getTotal(eventTypes);
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
    colors: ["#1A56DB", "#FDBA8C"],
    series: [
        {
            name: "Organic",
            color: "#1A56DB",
            data: [
                { x: "Mon", y: 231 },
                { x: "Tue", y: 122 },
                { x: "Wed", y: 63 },
                { x: "Thu", y: 421 },
                { x: "Fri", y: 122 },
                { x: "Sat", y: 323 },
                { x: "Sun", y: 111 },
            ],
        },
    ],
    chart: {
        type: "bar",
        height: "300px",
        fontFamily: "Inter, sans-serif",
        toolbar: {
            show: false,
        },
        animations: {
            enabled: false,
        },
    },
    plotOptions: {
        bar: {
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
            fontFamily: "Inter, sans-serif",
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
        enabled: false,
    },
    legend: {
        show: false,
    },
    xaxis: {
        floating: false,
        labels: {
            show: true,
            style: {
                fontFamily: "Inter, sans-serif",
                cssClass: 'text-xs font-normal fill-gray-500 dark:fill-gray-400'
            }
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
