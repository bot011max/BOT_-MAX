// Основные функции для веб-интерфейса

document.addEventListener('DOMContentLoaded', function() {
    console.log('Медицинский бот загружен');
    
    // Проверяем авторизацию
    checkAuth();
});

function checkAuth() {
    const token = localStorage.getItem('token');
    if (token) {
        // Обновляем интерфейс для авторизованного пользователя
        updateUIForAuth();
    }
}

function updateUIForAuth() {
    // Скрываем кнопки входа, показываем выход
    const loginBtns = document.querySelectorAll('.login-btn');
    loginBtns.forEach(btn => {
        btn.style.display = 'none';
    });
    
    // Добавляем кнопку выхода
    const logoutBtn = document.createElement('button');
    logoutBtn.className = 'btn logout-btn';
    logoutBtn.textContent = 'Выйти';
    logoutBtn.onclick = logout;
    document.querySelector('.container').appendChild(logoutBtn);
}

function logout() {
    localStorage.removeItem('token');
    window.location.reload();
}

// Функции для работы с API
async function apiRequest(url, method, data) {
    const token = localStorage.getItem('token');
    
    const response = await fetch(url, {
        method: method,
        headers: {
            'Content-Type': 'application/json',
            'Authorization': token ? `Bearer ${token}` : ''
        },
        body: data ? JSON.stringify(data) : undefined
    });
    
    return await response.json();
}
