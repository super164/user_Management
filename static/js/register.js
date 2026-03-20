async function submitRegister(event) {
    event.preventDefault();

    const form = document.getElementById('registerForm');
    const formData = new FormData(form);

    try {
        const response = await fetch('/register', {
            method: 'POST',
            headers: {
                'Accept': 'application/json'
            },
            body: formData
        });

        const text = await response.text();

        let result;
        try {
            result = JSON.parse(text);
        } catch (e) {
            alert(text || '注册失败');
            return;
        }

        if (result && result.success) {
            alert(result.message || '注册成功');
            window.location.href = '/login';
            return;
        }

        alert((result && result.message) ? result.message : '注册失败');
    } catch (error) {
        alert('网络请求失败，请稍后重试');
    }
}