// Classes
class DatePeriod {
    start
    end

    constructor(startDate, endDate) {
        this.start = startDate
        this.end = endDate
    }

    getRequestParamsInISOFormat() {
        return {
            start: this.start.toISOString(),
            end: this.end.toISOString()
        }
    }
}

class DataItemsList {
    items = []

    constructor(data = []) {
        this.items = data
    }
}

//Date intervals const's
const secondDuration = 1000
const dayDuration = 24 * 60 * 60 * secondDuration

const createDatePeriodToNow = (beforeDuration) => new DatePeriod(new Date(Date.now() - beforeDuration), new Date())
const dayDateInterval = () => createDatePeriodToNow(dayDuration)
const weekDateInterval = () => createDatePeriodToNow(7 * dayDuration)

// Endpoint const's
const serviceUrl = new URL("http://localhost:8080")
const dataEndpoint = "data"

// Chart settings
const chartDisplayFormats = {
    minute: 'HH:mm',
    hour: 'ddd/H'
}

function chartDefaultOptions() {
    return {
        scales: {
            x: {
                type: 'time',
                time: {
                    tooltipFormat: 'YYYY-MM-DD HH:mm',
                    displayFormats: chartDisplayFormats
                }
            }
        },
        borderWidth: 1,
        responsive: true,
        animation: false,
        interaction: {
            intersect: false
        },
        elements: {
            point: {
                radius: 0
            }
        },
        plugins: {
            zoom: {
                pan: {
                    enabled: true,
                    mode: 'xy',
                },
                zoom: {
                    wheel: {
                        enabled: true,
                    },
                    pinch: {
                        enabled: true,
                    },
                    mode: 'xy',
                }
            }
        }
    }
}

const renderInterval = 5 * secondDuration
const parameters = [
    "SellingProductsTotal",
    "SellingProductsOnWB",
    "SellingQuantityOnWB",
    "SellingProductsOnMP",
    "SellingQuantityOnMP",
    "ProductsOnWBStores",
    "QuantityOnWBStores",
    "ProductsOnExStores",
    "QuantityOnExStores",
]

async function getDataListByDateRange(datePeriod) {
    const requestUrl = new URL(dataEndpoint, serviceUrl)
    requestUrl.search = new URLSearchParams(datePeriod.getRequestParamsInISOFormat())

    const response = await fetch(requestUrl)
    return await response.json()
}

function updateData(dataItem) {
    const metricData = dataItem.MetricData
    parameters.forEach((value) => {
        setElementTextById(value, metricData[value].toLocaleString())
    })
}

function setElementTextById(id, text) {
    document.getElementById(id).innerText = text
}

function createCharts(dataItems) {
    return {
        lineCharts: createLineCharts(dataItems),
    }
}

function createLineCharts(dataItems) {
    const dataSetsLocal = dataSets(dataItems)
    const res = {}
    parameters.forEach((param) => {
        const key = `line-chart-${param}`
        const ctx = document.getElementById(key).getContext('2d')
        res[param] = new Chart(ctx, {
            type: 'line',
            data: {
                datasets: [dataSetsLocal[param]]
            },
            options: chartDefaultOptions()
        })
    })
    return res
}

function updateCharts(chartsTypes, dataItems) {
    const dataSetsLocal = dataSets(dataItems)
    Object.keys(chartsTypes).forEach((chartType) => {
        Object.keys(chartsTypes[chartType]).forEach((chartName) => {
            const chart = chartsTypes[chartType][chartName]
            chart.data = {
                datasets: [dataSetsLocal[chartName]]
            }
            chart.update()
        })
    })
}

function dataSets(dataItems) {
    const colors = [
        'red',
        'blue',
        'green',
        'brown',
        'magenta',
        'black',
        'indigo',
        'mediumvioletred',
        'darkgoldenrod',
    ]

    const counter = {
        SellingProductsTotal: [],
        SellingProductsOnWB: [],
        SellingQuantityOnWB: [],
        SellingProductsOnMP: [],
        SellingQuantityOnMP: [],
        ProductsOnWBStores: [],
        QuantityOnWBStores: [],
        ProductsOnExStores: [],
        QuantityOnExStores: [],
    }

    dataItems.items.forEach((item) => {
        parameters.forEach((parameter) => {
            counter[parameter].push({
                x: item.Time,
                y: item.MetricData[parameter]
            })
        })
    })

    const sets = {}
    Object.keys(counter).forEach((key, index) => {
        sets[key] = {
            label: key,
            data: counter[key],
            borderColor: colors[index],
            backgroundColor: colors[index],
        }
    })
    return sets
}

function getDateInterval() {
    switch (currentDateIntervalOption) {
        case "week":
            return weekDateInterval()
        default:
            return dayDateInterval()
    }
}

function setDateIntervalScales() {
// chart.options.scales = getDateIntervalScales()
}

function render(data) {
    if (data.length < 1) return

    const lastData = data.at(-1)
    updateData(lastData)

    updateCharts(charts, new DataItemsList(data))
}

let currentDateIntervalOption = "day"
const charts = createCharts(new DataItemsList())

function changeOption(e) {
    currentDateIntervalOption = e.target.value
}

const selection = document.getElementById("intervalSelect");
selection.addEventListener("change", changeOption);

window.onload = async function () {
    const test = async () => render(await getDataListByDateRange(getDateInterval()))
    await test()

    setInterval(test, renderInterval)
};