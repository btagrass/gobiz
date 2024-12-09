package svc

import (
	"log/slog"
	"os"

	"github.com/btagrass/gobiz/svc"
	"github.com/samber/lo"
)

func Init(svcs ...string) {
	if len(svcs) == 0 || lo.Contains(svcs, "dept") {
		svc.Inject(NewDeptSvc)
		err := svc.Migrate("INSERT IGNORE INTO sys_dept VALUES (300000000000001, '2023-01-29 00:00:00.000', NULL, NULL, 0, 'Kskj', '', '', 1)")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}
	if len(svcs) == 0 || lo.Contains(svcs, "user") {
		svc.Inject(NewUserSvc)
		err := svc.Migrate("INSERT IGNORE INTO sys_user VALUES (300000000000001, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, 'admin', NULL, '15800000000', '$2a$10$enX7NxYTZZo9yLJQN6jXF.B6FGg7d9Q5eTW5off94hJZSa5AO9av2', 0)")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}
	if len(svcs) == 0 || lo.Contains(svcs, "dict") {
		svc.Inject(NewDictSvc)
		err := svc.Migrate(
			"INSERT IGNORE INTO sys_dict VALUES (300000001060101, '2023-01-29 00:00:00.000', NULL, NULL, 'CacheType', 1, 'Local', 1)",
			"INSERT IGNORE INTO sys_dict VALUES (300000001060102, '2023-01-29 00:00:00.000', NULL, NULL, 'CacheType', 2, 'Remote', 2)",
		)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}
	if len(svcs) == 0 || lo.Contains(svcs, "job") {
		svc.Inject(NewJobSvc)
		err := svc.Migrate(
			"INSERT IGNORE INTO sys_dict VALUES (300000001070101, '2023-01-29 00:00:00.000', NULL, NULL, 'JobState', 1, 'Started', 1)",
			"INSERT IGNORE INTO sys_dict VALUES (300000001070102, '2023-01-29 00:00:00.000', NULL, NULL, 'JobState', 0, 'Stopped', 2)",
		)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}
	if len(svcs) == 0 || lo.Contains(svcs, "resource") {
		svc.Inject(NewResourceSvc)
		err := svc.Migrate(
			"INSERT IGNORE INTO sys_resource VALUES (300000000000001, '2023-01-29 00:00:00.000', NULL, NULL, 0, 'SystemSettings', 1, 'Setting', '/sys', NULL, 100)",
			"INSERT IGNORE INTO sys_resource VALUES (300000000000101, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, 'DepartmentManagement', 1, 'OfficeBuilding', '/sys/depts', NULL, 1)",
			"INSERT IGNORE INTO sys_resource VALUES (300000000000102, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, 'UserManagement', 1, 'User', '/sys/users', NULL, 2)",
			"INSERT IGNORE INTO sys_resource VALUES (300000000000103, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, 'RoleManagement', 1, 'Avatar', '/sys/roles', NULL, 3)",
			"INSERT IGNORE INTO sys_resource VALUES (300000000000104, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, 'ResourceManagement', 1, 'Menu', '/sys/resources', NULL, 4)",
			"INSERT IGNORE INTO sys_resource VALUES (300000000000105, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, 'DictionaryManagement', 1, 'Files', '/sys/dicts', NULL, 5)",
			"INSERT IGNORE INTO sys_resource VALUES (300000000000106, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, 'CacheManagement', 1, 'Coin', '/sys/caches', NULL, 6)",
			"INSERT IGNORE INTO sys_resource VALUES (300000000000107, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, 'JobManagement', 1, 'Clock', '/sys/jobs', NULL, 7)",
			"INSERT IGNORE INTO sys_resource VALUES (300000000000108, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, 'VisitRecords', 1, 'PieChart', '/sys/visits', NULL, 8)",
			"INSERT IGNORE INTO sys_resource VALUES (300000000000109, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, 'ServerInformation', 1, 'Histogram', '/sys/servers', NULL, 9)",
			"INSERT IGNORE INTO sys_dict VALUES (300000001040101, '2023-01-29 00:00:00.000', NULL, NULL, 'ResourceType', 1, 'Menu', 1)",
			"INSERT IGNORE INTO sys_dict VALUES (300000001040102, '2023-01-29 00:00:00.000', NULL, NULL, 'ResourceType', 2, 'Authority', 2)",
			"INSERT IGNORE INTO sys_dict VALUES (300000001040201, '2023-01-29 00:00:00.000', NULL, NULL, 'ResourceAct', 1, 'GET', 1)",
			"INSERT IGNORE INTO sys_dict VALUES (300000001040202, '2023-01-29 00:00:00.000', NULL, NULL, 'ResourceAct', 2, 'DELETE', 2)",
			"INSERT IGNORE INTO sys_dict VALUES (300000001040203, '2023-01-29 00:00:00.000', NULL, NULL, 'ResourceAct', 3, 'POST', 3)",
		)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}
	if len(svcs) == 0 || lo.Contains(svcs, "role") {
		svc.Inject(NewRoleSvc)
		err := svc.Migrate("INSERT IGNORE INTO sys_role VALUES (300000000000001, '2023-01-29 00:00:00.000', NULL, NULL, 'Admin')")
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}
	if len(svcs) == 0 || lo.Contains(svcs, "server") {
		svc.Inject(NewServerSvc)
	}
	if len(svcs) == 0 || lo.Contains(svcs, "upgrade") {
		svc.Inject(NewUpgradeSvc)
	}
	if len(svcs) == 0 || lo.Contains(svcs, "visit") {
		svc.Inject(NewVisitSvc)
	}
}
