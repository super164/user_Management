// 定义全局变量，由 HTML 在引入此文件前定义
// const currentUserRole = ...
// const currentUserId = ...

function handleUsersPageNotice() {
    const urlParams = new URLSearchParams(window.location.search);
    const msg = urlParams.get('msg');
    const open = urlParams.get('open');
    const flashMsg = getCookie('flash_msg');
    const flashOpen = getCookie('flash_open');

    if (flashOpen === 'create' || open === 'create') {
        showCreateModal();
    }

    if (flashMsg) {
        alert(decodeURIComponent(flashMsg.replace(/\+/g, '%20')));
        deleteCookie('flash_msg');
    }

    if (msg) {
        alert(msg);
    }

    if (flashOpen) {
        deleteCookie('flash_open');
    }

    if (msg || open) {
        const newUrl = window.location.protocol + "//" + window.location.host + window.location.pathname;
        window.history.replaceState({ path: newUrl }, '', newUrl);
    }
}

function getCookie(name) {
    const prefix = name + '=';
    const parts = document.cookie.split('; ');
    for (const part of parts) {
        if (part.startsWith(prefix)) {
            return part.slice(prefix.length);
        }
    }
    return null;
}

function deleteCookie(name) {
    document.cookie = name + '=; Path=/; Max-Age=0';
}

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', handleUsersPageNotice);
} else {
    handleUsersPageNotice();
}

function openEditModal(user) {
    const modal = document.getElementById('editUserModal');

    // 1. 填充数据
    document.getElementById('modal_user_id').value = user.id;
    document.getElementById('modal_username').value = user.username;
    // 清空密码框
    const passwordInput = document.querySelector('input[name="password"]');
    if (passwordInput) passwordInput.value = '';

    // 头像预览
    const avatarSrc = user.avatar ? '/uploads/' + user.avatar : 'https://via.placeholder.com/40';
    document.getElementById('modal_avatar_preview').src = avatarSrc;

    // 2. 权限控制
    const statusSelect = document.getElementById('modal_status');
    const usernameInput = document.getElementById('modal_username');

    // 如果不是管理员，禁用状态选择
    if (currentUserRole !== 'admin') {
        if (statusSelect) statusSelect.disabled = true;
        // 如果不是改自己，甚至可能禁用用户名修改等
        if (user.id !== currentUserId) {
            if (usernameInput) usernameInput.readOnly = true;
        }
    } else {
        if (statusSelect) statusSelect.disabled = false;
        if (usernameInput) usernameInput.readOnly = false;
    }

    // 3. 显示模态框
    modal.classList.remove('hidden');
}

function closeEditModal() {
    document.getElementById('editUserModal').classList.add('hidden');
}

// 点击遮罩层关闭编辑模态框
const editUserModal = document.getElementById('editUserModal');
if (editUserModal) {
    editUserModal.addEventListener('click', function(e) {
        if (e.target === this) {
            closeEditModal();
        }
    });
}

// 打开新建用户模态框并重置表单
function openCreateModal() {
    const modal = document.getElementById('createUserModal');
    const form = document.getElementById('createUserForm');
    
    if (form) {
        form.reset();
    }
    
    if (modal) {
        modal.classList.remove('hidden');
    }
}

function showCreateModal() {
    const modal = document.getElementById('createUserModal');
    if (modal) {
        modal.classList.remove('hidden');
    }
}

// 点击遮罩层关闭新建用户模态框
const createUserModal = document.getElementById('createUserModal');
if (createUserModal) {
    createUserModal.addEventListener('click', function(e) {
        if (e.target === this) {
            createUserModal.classList.add('hidden');
        }
    });
}

// 提交更新表单
async function submitUpdateForm(event) {
    event.preventDefault(); // 阻止表单默认跳转提交

    const form = document.getElementById('updateUserForm');
    const formData = new FormData(form);

    try {
        const response = await fetch('/users/update', {
            method: 'POST',
            body: formData
        });

        // 处理被重定向到登录页的情况（如会话过期或被踢出）
        if (response.redirected && response.url.includes('/login')) {
            window.location.href = response.url;
            return;
        }

        // 处理非 200 响应
        if (!response.ok) {
            const errorText = await response.text();
            try {
                const errorJson = JSON.parse(errorText);
                alert(errorJson.message || '更新失败');
            } catch (e) {
                alert('更新失败: ' + errorText);
            }
            return;
        }

        // 处理成功响应
        const result = await response.json();
        if (result.success) {
            alert('更新成功！');
            window.location.reload();
        } else {
            alert(result.message || '更新失败');
        }

    } catch (error) {
        console.error('Error:', error);
        alert('网络请求失败，请稍后重试');
    }
}

async function submitCreateUser(event) {
    if (event && typeof event.preventDefault === 'function') {
        event.preventDefault();
    }

    const form = document.getElementById('createUserForm');
    const formData = new FormData(form);

    try {
        const response = await fetch('/users/create', {
            method: 'POST',
            headers: {
                'Accept': 'application/json'
            },
            body: formData
        });

        if (response.redirected && response.url.includes('/login')) {
            window.location.href = response.url;
            return;
        }

        const text = await response.text();
        let result;
        try {
            result = JSON.parse(text);
        } catch (e) {
            alert(text || '创建失败');
            showCreateModal();
            return;
        }

        if (result && result.success) {
            alert(result.message || '创建成功');
            window.location.reload();
            return;
        }

        alert((result && result.message) ? result.message : '创建失败');
        showCreateModal();
    } catch (error) {
        alert('网络请求失败，请稍后重试');
        showCreateModal();
    }

    return false;
}

// 全局功能
function setPageTitle(title) {
    const pageTitle = document.getElementById('pageTitle');
    if (pageTitle) {
        pageTitle.textContent = title;
    }
    document.title = title + ' - 后台管理系统';
}

function logout() {
    if (confirm('确定要退出系统吗？')) {
        localStorage.clear();
        window.location.href = '/logout';
    }
}

// 初始化图标
if (typeof feather !== 'undefined') {
    feather.replace();
}
