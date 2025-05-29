# imd-seat-be
信息管理学院选座系统后端

## 目录结构说明
```
.
├── api/            # API接口定义文件
│   ├── index.api   # 代码生成入口API
│   └── *.api       # 各模块API定义
├── etc/            # 配置文件
│   └── config.yaml # 应用配置
├── internal/
│   ├── config/     # 配置结构体定义
│   ├── handler/    # 路由处理层
│   ├── logic/      # 业务逻辑层
│   ├── model/      # 数据模型层
│   ├── svc/        # 服务上下文
│   ├── types/      # 请求/响应DTO
│   └── sql/        # 数据表设计
├── go.mod          # Go模块定义
├── go.sum          # 依赖校验
├── main.go         # 主入口
└── Makefile        # 构建脚本
```

## API开发流程
1. 在对应的`.api`文件中定义接口(如`reservation.api`)
```go
@server(
    prefix: /api/reservation
    group: reservation
)
service imd-seat-be {
    @handler ReserveSeatHandler
    post /reserve (ReserveRequest) returns (ReserveResponse)
    
    @handler CancelReservationHandler  
    post /cancel (CancelRequest) returns (CancelResponse)
}

type (
    ReserveRequest {
        SeatID int64 `json:"seatId"`
        Date   string `json:"date"`
    }
    
    ReserveResponse {
        ReservationID int64 `json:"reservationId"`
    }
    
    CancelRequest {
        ReservationID int64 `json:"reservationId"`
    }
    
    CancelResponse {
        Success bool `json:"success"`
    }
)
```

2. 生成代码
```bash
make api
```

3. 实现业务逻辑
```go
// internal/logic/reserveseatlogic.go
func (l *ReserveSeatLogic) ReserveSeat(req *types.ReserveRequest) (*types.ReserveResponse, error) {
    // 1. 参数校验
    // 2. 检查座位可用性
    // 3. 创建预约记录
    return &types.ReserveResponse{ReservationID: newID}, nil
}
```

## 数据库设计
主要表结构见`internal/sql/data.sql`，包含：
- 座位表(seats)  
- 预约表(reservations)
- 房间表(rooms)
- 用户表(users)

## 错误码定义

### 用户模块状态码（1000-1999）
| 错误码 | 描述 |
|--------|------|
| 1001 | 学号或密码错误 |
| 1002 | 用户不存在 |
| 1003 | 身份验证失败 |
| 1004 | 鉴权失败 |

### 爬虫模块状态码（2000-2999）
| 错误码 | 描述 |
|--------|------|
| 2001 | 爬虫错误 |
| 2002 | ccnu服务器错误 |

### 数据库模块状态码（3000-3999）
| 错误码 | 描述 |
|--------|------|
| 3001 | 数据库查询失败 |
| 3002 | 数据库创建失败 |
| 3003 | 数据库更新失败/删除失败 |

### 预约模块状态码（4000-4999）
| 错误码 | 描述 |
|--------|------|
| 4001 | 预约请求不合规 |
| 4002 | 已经签到过了 |
| 4003 | 您还未签到 |

### 默认错误码
| 错误码 | 描述 |
|--------|------|
| 5000 | 非预设错误 |

### 预约状态码
| 状态码 | 描述 |
|--------|------|
| booked | 已预约 |
| effective | 已生效未签到 |
| checked in | 已签到 |
| cancelled | 已取消 |
| completed | 预约已完成 |
| one default | 已未到一次 |
| violated | 已违约 |
| available | 可预约 |

## 接口设计
[接口设计](https://apifox.com/apidoc/shared-2723d5c5-2c20-467d-a3d6-5138f50b0e4b)
