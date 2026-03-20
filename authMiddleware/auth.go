package authMiddleware

import (
	"net/http"
	"userManagement/dao"
	"userManagement/model"
	"userManagement/session"
)

// AuthMiddleware 登录验证包装器
func AuthMiddleware(next func(http.ResponseWriter, *http.Request, model.User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 检查Session是否存在
		currentUser, ok := session.GetSession(r)
		if !ok {
			// 未登录，跳转到登录页
			http.Redirect(w, r, "/login", 302)
			return
		}
		//查询数据库确认用户最新状态
		dbUser, err := dao.GetUserByID(currentUser.ID)
		if err != nil {
			// 数据库查询出错，安全起见认为是认证失败或系统错误
			http.Error(w, "系统错误", http.StatusInternalServerError)
			return
		}
		// 3. 检查用户是否存在 (可能已被删除)
		if dbUser == nil {
			// 用户不存在，销毁 Session 并踢出
			session.DestroySession(w, r)
			http.Redirect(w, r, "/login?msg=account_deleted", 302)
			return
		}

		// 4. 检查用户状态 (可能已被封禁)
		// 假设 status=0 代表禁用/封禁
		if dbUser.Status == 0 {
			// 账号已封禁，销毁 Session 并踢出
			session.DestroySession(w, r)
			// 强制跳转回登录页，带上错误信息
			http.Redirect(w, r, "/login?msg=account_banned", 302)
			return
		}
		// 验证通过，将 currentUser 传递给业务逻辑
		next(w, r, currentUser)
	}
}
