const lineOptions = {
    chart: {
        height: "100%",
        maxWidth: "100%",
        type: "area",
        dropShadow: {
            enabled: false,
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
            show: false,
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
            data: [6500, 6418, 6456, 6526, 6356, 6456],
            color: "#1A56DB",
        },
    ],
    xaxis: {
        categories: ['01 February', '02 February', '03 February', '04 February', '05 February', '06 February', '07 February'],
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
    series: [52.8, 26.8, 20.4],
    colors: ["#1C64F2", "#16BDCA", "#9061F9"],
    chart: {
        height: 420,
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
    labels: ["Direct", "Organic search", "Referrals"],
    dataLabels: {
        enabled: true,
        style: {
            fontFamily: "Space Grotesk",
        },
    },
    legend: {
        position: "bottom",
        fontFamily: "Inter, sans-serif",
    },
    yaxis: {
        labels: {
            formatter: function(value) {
                return value + "%"
            },
        },
    },
    xaxis: {
        labels: {
            formatter: function(value) {
                return value + "%"
            },
        },
        axisTicks: {
            show: false,
        },
        axisBorder: {
            show: false,
        },
    },
}

const barOptions = {
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
        height: "320px",
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

if (document.getElementById("area-chart") && typeof ApexCharts !== 'undefined') {
    const chart = new ApexCharts(document.getElementById("area-chart"), lineOptions);
    chart.render();
}

if (document.getElementById("pie-chart") && typeof ApexCharts !== 'undefined') {
    const chart = new ApexCharts(document.getElementById("pie-chart"), pieOptions);
    chart.render();
}

if (document.getElementById("column-chart") && typeof ApexCharts !== 'undefined') {
    const chart = new ApexCharts(document.getElementById("column-chart"), barOptions);
    chart.render();
}
