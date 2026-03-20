async function submitLogin(event) {
    event.preventDefault(); // 阻止表单默认提交

    const form = document.getElementById('loginForm');
    const formData = new FormData(form);

    try {
        const response = await fetch('/login', {
            method: 'POST',
            body: formData
        });

        if (response.redirected) {
            window.location.href = response.url;
            return;
        }

        const text = await response.text();

        // 简单的错误文本判断
        if (text.includes("失败") || text.includes("错误") || text.includes("禁用") || text.includes("不存在")) {
            alert(text);
        } else {
            // 尝试解析是否是 JSON 格式的错误
            try {
                const json = JSON.parse(text);
                if (json.success) {
                    window.location.href = '/users';
                } else {
                    alert(json.message || '登录失败');
                }
            } catch (e) {
                // 如果不是 JSON 也不是错误文本，可能是页面 HTML，直接替换显示
                document.documentElement.innerHTML = text;
            }
        }

    } catch (error) {
        console.error('Error:', error);
        alert('网络请求失败，请稍后重试');
    }
}

// 页面加载时检查 URL 参数中的错误信息
window.onload = function() {
    const urlParams = new URLSearchParams(window.location.search);
    const msg = urlParams.get('msg');
    if (msg === 'account_deleted') {
        alert('您的账号已被删除，无法登录，请联系管理员。');
    } else if (msg === 'account_banned') {
        alert('您的账号已被禁用，无法登录，请联系管理员。');
    }
    // 清除 URL 参数
    if (msg) {
        const newUrl = window.location.protocol + "//" + window.location.host + window.location.pathname;
        window.history.replaceState({
            path: newUrl
        }, '', newUrl);
    }
}