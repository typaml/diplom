document.addEventListener('DOMContentLoaded', async function() {
    let callsChart;
    let salesChart; // Объявляем переменную salesChart здесь

    async function fetchData(url) {
        const response = await fetch(url);
        return response.json();
    }

    async function initCharts(period) {
        let url;
        switch (period) {
            case 'day':
                url = '/calls?period=day';
                break;
            case 'week':
                url = '/calls?period=week';
                break;
            case 'month':
                url = '/calls?period=month';
                break;
            case 'year':
                url = '/calls?period=year';
                break;
            default:
                url = '/calls';
        }

        const callsData = await fetchData(url);
        const chartOptions = {
            scales: {
                y: {
                    beginAtZero: true
                }
            }
        };

        // Уничтожаем существующий график, если он существует
        if (callsChart) {
            callsChart.destroy();
        }

        const callsChartCtx = document.getElementById('calls-chart').getContext('2d');
        callsChart = new Chart(callsChartCtx, {
            type: 'line',
            data: callsData,
            options: chartOptions
        });
    }

    async function initSalesChart(period) {
        let url;
        switch (period) {
            case 'day':
                url = '/sales?period=day';
                break;
            case 'week':
                url = '/sales?period=week';
                break;
            case 'month':
                url = '/sales?period=month';
                break;
            case 'year':
                url = '/sales?period=year';
                break;
            default:
                url = '/sales';
        }


        const salesData = await fetchData(url);


        const chartOptions = {
            scales: {
                y: {
                    beginAtZero: true
                }
            }
        };

        // Уничтожаем существующий график, если он существует
        if (salesChart) {
            salesChart.destroy();
        }

        const salesChartCtx = document.getElementById('sales-chart').getContext('2d');
        salesChart = new Chart(salesChartCtx, {
            type: 'bar',
            data: salesData,
            options: chartOptions
        });
    }

    await initSalesChart('day');

    await initCharts('day'); // Инициализируем график для начального периода (день)

    const periodSelect = document.getElementById('period-select');
    periodSelect.addEventListener('change', function() {
        const selectedPeriod = this.value;
        initCharts(selectedPeriod); // Обновляем данные для графика при изменении периода
        initSalesChart(selectedPeriod);
    });
});
