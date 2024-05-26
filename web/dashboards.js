document.addEventListener('DOMContentLoaded', function() {
    // Запрос данных для графика "Клиенты по статусам"
    fetch('/getclientsstatus')
        .then(response => response.json())
        .then(data => {
            const ctx = document.getElementById('statusChart').getContext('2d');
            const labels = data.map(item => item.status);
            const counts = data.map(item => item.count);

            const backgroundColors = brightColors(data.length); // Генерация ярких цветов
            const borderColors = backgroundColors.map(color => color.replace('0.2', '1')); // Установка альфа-канала в 1 для border

            new Chart(ctx, {
                type: 'bar',
                data: {
                    labels: labels,
                    datasets: [{
                        label: 'Количество клиентов',
                        data: counts,
                        backgroundColor: backgroundColors,
                        borderColor: borderColors,
                        borderWidth: 1
                    }]
                },
                options: {
                    plugins: {
                        title: {
                            display: true,
                            text: 'Клиенты по статусам',
                            font: {
                                size: 18
                            }
                        },
                        legend: {
                            display: false
                        }
                    },
                    scales: {
                        y: {
                            beginAtZero: true,
                            title: {
                                display: true,
                                text: 'Количество клиентов'
                            },
                            grid: {
                                borderDash: [5, 5]
                            }
                        },
                        x: {
                            title: {
                                display: true,
                                text: 'Статус клиента'
                            },
                            grid: {
                                borderDash: [5, 5]
                            }
                        }
                    },
                    layout: {
                        padding: 5
                    },
                    elements: {
                        bar: {
                            borderRadius: 5,
                            borderSkipped: 'bottom'
                        }
                    }
                }
            });
        })
        .catch(error => console.error('Error:', error));

    // Запрос данных для графика "Клиенты по маркетинговым рассылкам"
    fetch('/getclientsevents')
        .then(response => response.json())
        .then(data => {
            const ctx = document.getElementById('marketingChart').getContext('2d');
            const labels = data.map(item => item.marketing_event);
            const counts = data.map(item => item.count);

            new Chart(ctx, {
                type: 'pie', // Измененный тип графика
                data: {
                    labels: labels,
                    datasets: [{
                        label: 'Количество клиентов',
                        data: counts,
                        backgroundColor: brightColors(data.length), // Генерация ярких цветов
                        borderColor: '#fff',
                        borderWidth: 1
                    }]
                },
                options: {
                    plugins: {
                        title: {
                            display: true,
                            text: 'Клиенты по маркетинговым ивентам',
                            font: {
                                size: 18
                            }
                        },
                        legend: {
                            display: true,
                            position: 'right' // Положение легенды
                        }
                    },
                    layout: {
                        padding: 10
                    }
                }
            });
        })
        .catch(error => console.error('Error:', error));

    // Запрос данных для графика "Клиенты по регионам"
    fetch('/getclientsregion')
    .then(response => response.json())
    .then(data => {
        const ctx = document.getElementById('regionChart').getContext('2d');
        const labels = data.map(item => item.region_name);
        const counts = data.map(item => item.count);

        const backgroundColors = brightColors(data.length); // Генерация ярких цветов
        const borderColors = backgroundColors.map(color => color.replace('0.2', '1')); // Установка альфа-канала в 1 для border

        new Chart(ctx, {
            type: 'doughnut', // Изменение типа графика на "донат"
            data: {
                labels: labels,
                datasets: [{
                    label: 'Количество клиентов',
                    data: counts,
                    backgroundColor: backgroundColors,
                    borderColor: borderColors,
                    borderWidth: 1
                }]
            },
            options: {
                plugins: {
                    title: {
                        display: true,
                        text: 'Клиенты по регионам',
                        font: {
                            size: 18
                        }
                    },
                    legend: {
                        display: true,
                        position: 'right' // Положение легенды
                    }
                },
                layout: {
                    padding: 10
                }
            }
        });
    })
    .catch(error => console.error('Error:', error));
    fetch('/getcallstoday')
    .then(response => response.json())
    .then(data => {
        const ctx = document.getElementById('callsChart').getContext('2d');
        const label = 'Сегодня'; // Метка для сегодняшнего дня
        const count = data.count; // Полученное количество звонков за сегодня

        const backgroundColor = randomBrightColor(); // Генерация яркого цвета

        new Chart(ctx, {
            type: 'bar',
            data: {
                labels: [label],
                datasets: [{
                    label: 'Количество звонков',
                    data: [count],
                    backgroundColor: backgroundColor,
                    borderColor: backgroundColor.replace('0.2', '1'), // Установка альфа-канала в 1 для border
                    borderWidth: 1
                }]
            },
            options: {
                plugins: {
                    title: {
                        display: true,
                        text: 'Количество звонков за сегодня',
                        font: {
                            size: 18
                        }
                    },
                    legend: {
                        display: false
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        title: {
                            display: true,
                            text: 'Количество звонков'
                        },
                        grid: {
                            borderDash: [5, 5]
                        }
                    },
                    x: {
                        title: {
                            display: true,
                            text: 'День'
                        },
                        grid: {
                            borderDash: [5, 5]
                        }
                    }
                },
                layout: {
                    padding: 20
                },
                elements: {
                    bar: {
                        borderRadius: 5,
                        borderSkipped: 'bottom'
                    }
                }
            }
        });
    })
    .catch(error => console.error('Error:', error));
    fetch('/getclientsbyuser')
    .then(response => response.json())
    .then(data => {
        const ctx = document.getElementById('userClientsChart').getContext('2d');
        const labels = data.map(item => item.UserName);
        const counts = data.map(item => item.ClientCount);
        const backgroundColors = brightColors(data.length); // Генерация ярких цветов для каждого столбца

        new Chart(ctx, {
            type: 'bar',
            data: {
                labels: labels,
                datasets: [{
                    label: 'Количество клиентов',
                    data: counts,
                    backgroundColor: backgroundColors, // Использование сгенерированных цветов
                    borderColor: backgroundColors.map(color => color.replace('0.2', '1')), // Обводка того же цвета
                    borderWidth: 1
                }]
            },
            options: {
                plugins: {
                    title: {
                        display: true,
                        text: 'Количество клиентов по ответственным',
                        font: {
                            size: 18
                        }
                    },
                    legend: {
                        display: false
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        title: {
                            display: true,
                            text: 'Количество клиентов'
                        },
                        grid: {
                            borderDash: [5, 5]
                        }
                    },
                    x: {
                        title: {
                            display: true,
                            text: 'Ответственный'
                        },
                        grid: {
                            borderDash: [5, 5]
                        }
                    }
                },
                layout: {
                    padding: 20
                },
                elements: {
                    bar: {
                        borderRadius: 5,
                        borderSkipped: 'bottom'
                    }
                }
            }
        });
    })
    .catch(error => console.error('Error:', error));
});

// Функция для генерации ярких и насыщенных цветов
function brightColors(count) {
    const colors = [];
    for (let i = 0; i < count; i++) {
        colors.push(randomBrightColor());
    }
    return colors;
}

// Функция для генерации случайного яркого цвета
function randomBrightColor() {
    const r = Math.floor(Math.random() * 256); // Красный
    const g = Math.floor(Math.random() * 256); // Зеленый
    const b = Math.floor(Math.random() * 256); // Синий
    return `rgb(${r}, ${g}, ${b})`;
}
function redirectToFilteredPage(event, filter) {
    event.preventDefault(); // Предотвращаем переход по ссылке

    // Перенаправляем пользователя на вкладку с нужным фильтром
    window.location.href = '/clients?filterStatus=' + filter;
}

const createClientBtn = document.getElementById('createClientBtn');
    const createClientModal = document.getElementById('createClientModal');
    const closeModal = document.getElementsByClassName('close')[0];
    const createClientForm = document.getElementById('createClientForm');

    createClientBtn.onclick = function() {
        createClientModal.style.display = 'block';
    }

    closeModal.onclick = function() {
        createClientModal.style.display = 'none';
    }

    window.onclick = function(event) {
        if (event.target == createClientModal) {
            createClientModal.style.display = 'none';
        }
    }

    createClientForm.addEventListener('submit', function(event) {
        event.preventDefault(); // Предотвращаем стандартное поведение формы

        const formData = new FormData(createClientForm); // Получаем данные формы
        const jsonData = {};
        formData.forEach((value, key) => {
            jsonData[key] = value;
        });

        fetch('/createclient', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(jsonData), // Преобразуем данные в JSON и отправляем на сервер
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to create client');
            }
            return response.json();
        })
        .then(data => {
            console.log('Client created successfully:', data);
            // Дополнительные действия при успешном создании клиента
        })
        .catch(error => {
            console.error('Error:', error);
            // Обработка ошибок при создании клиента
        })
        .finally(() => {
            createClientModal.style.display = 'none'; // Закрываем модальное окно после отправки данных
        });
        
    });
    