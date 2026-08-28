package role

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"iot-zero/services/platform-api/internal/svc"
	"iot-zero/services/platform-api/model"
)

func newMenuTreeLogicWithMock(t *testing.T, ctx context.Context) (*GetMenuTreeLogic, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	conn := sqlx.NewSqlConnFromDB(db)
	return NewGetMenuTreeLogic(ctx, &svc.ServiceContext{
		RoleModel:     model.NewRoleModel(conn),
		MenuModel:     model.NewMenuModel(conn),
		RoleMenuModel: model.NewRoleMenuModel(conn),
	}), mock
}

func TestGetMenuTreeLogic_GetMenuTree(t *testing.T) {
	now := time.Now()

	t.Run("unauthorized", func(t *testing.T) {
		ast := assert.New(t)
		l, _ := newMenuTreeLogicWithMock(t, context.Background())
		_, err := l.GetMenuTree()
		ast.ErrorIs(err, errUnauthorized)
	})

	t.Run("role not found", func(t *testing.T) {
		ast := assert.New(t)
		ctx := context.WithValue(context.Background(), "roleType", "missing")
		l, mock := newMenuTreeLogicWithMock(t, ctx)

		mock.ExpectQuery("select .+ from `role` where `role_type` = \\? and `deleted_at` is null").
			WithArgs("missing").
			WillReturnError(sqlx.ErrNotFound)

		_, err := l.GetMenuTree()
		ast.ErrorIs(err, errRoleNotFound)
	})

	t.Run("empty menus", func(t *testing.T) {
		ast := assert.New(t)
		ctx := context.WithValue(context.Background(), "roleType", "admin")
		l, mock := newMenuTreeLogicWithMock(t, ctx)

		mock.ExpectQuery("select .+ from `role` where `role_type` = \\? and `deleted_at` is null").
			WithArgs("admin").
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "role_type", "role_name", "enable", "tenant_id"}).
				AddRow(uint64(2), now, now, nil, "admin", "Admin", "Enable", uint64(1)))
		mock.ExpectQuery("select `menu_id` from `role_menu` where `role_id` = \\?").
			WithArgs(uint64(2)).
			WillReturnRows(sqlmock.NewRows([]string{"menu_id"}))

		resp, err := l.GetMenuTree()
		ast.NoError(err)
		ast.Empty(resp.List)
	})

	t.Run("builds tree through role_menu and ancestors", func(t *testing.T) {
		ast := assert.New(t)
		ctx := context.WithValue(context.Background(), "roleType", "admin")
		l, mock := newMenuTreeLogicWithMock(t, ctx)

		mock.ExpectQuery("select .+ from `role` where `role_type` = \\? and `deleted_at` is null").
			WithArgs("admin").
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "role_type", "role_name", "enable", "tenant_id"}).
				AddRow(uint64(2), now, now, nil, "admin", "Admin", "Enable", uint64(1)))
		mock.ExpectQuery("select `menu_id` from `role_menu` where `role_id` = \\?").
			WithArgs(uint64(2)).
			WillReturnRows(sqlmock.NewRows([]string{"menu_id"}).AddRow(uint64(2)).AddRow(uint64(54)))
		mock.ExpectQuery("select .+ from `menu` where `id` in \\(\\?,\\?\\) order by `id`").
			WithArgs(uint64(2), uint64(54)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "parent_id", "menu_key", "title", "has_child"}).
				AddRow(uint64(2), uint64(1), "user", "用户管理", "N").
				AddRow(uint64(54), uint64(0), "videoCenter", "视频中心", "N"))
		mock.ExpectQuery("select .+ from `menu` where `id` in \\(\\?\\) order by `id`").
			WithArgs(uint64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "parent_id", "menu_key", "title", "has_child"}).
				AddRow(uint64(1), uint64(0), "system", "系统管理", "Y"))

		resp, err := l.GetMenuTree()
		ast.NoError(err)
		ast.Len(resp.List, 2)

		byKey := make(map[string]int)
		for i, item := range resp.List {
			byKey[item.MenuKey] = i
		}
		ast.Contains(byKey, "system")
		ast.Contains(byKey, "videoCenter")

		system := resp.List[byKey["system"]]
		ast.Equal(uint64(1), system.Id)
		ast.Equal("Y", system.HasChild)
		ast.Len(system.Children, 1)
		ast.Equal("user", system.Children[0].MenuKey)
		ast.Equal(uint64(1), system.Children[0].ParentId)

		video := resp.List[byKey["videoCenter"]]
		ast.Equal(uint64(54), video.Id)
		ast.Empty(video.Children)
	})
}

func TestBuildMenuTree(t *testing.T) {
	ast := assert.New(t)
	tree := buildMenuTree([]*model.Menu{
		{Id: 2, ParentId: 1, MenuKey: "user", Title: "用户管理", HasChild: "N"},
		{Id: 1, ParentId: 0, MenuKey: "system", Title: "系统管理", HasChild: "Y"},
		{Id: 3, ParentId: 1, MenuKey: "role", Title: "角色管理", HasChild: "N"},
	})
	ast.Len(tree, 1)
	ast.Equal("system", tree[0].MenuKey)
	ast.Len(tree[0].Children, 2)
	ast.Equal("user", tree[0].Children[0].MenuKey)
	ast.Equal("role", tree[0].Children[1].MenuKey)
}
