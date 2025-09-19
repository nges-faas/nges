package queryparser

import "strings"

// type & constant definition

// Action 操作
type Action int

const (
	Query  Action = 1
	Insert Action = 2
	Update Action = 4
	Delete Action = 8
)

// Is 判断是否为指定操作
func (a Action) Is(p Action) bool {
	return a == p
}

//func (a Action) Has(p Action) bool {
//	return a&p > 0
//}

// Connector 连接符
type Connector int

const (
	ConnectorAnd Connector = 0
	ConnectorOr  Connector = 1
	ConnectorNot Connector = 2
)

// Is 判断是否为指定连接符
func (c Connector) Is(p Connector) bool {
	return c == p
}

// ValueType 值类型
type ValueType int

const (
	ValueTypeUnknown  ValueType = iota // any
	ValueTypeNull                      // nil
	ValueTypeText                      // string
	ValueTypeInt                       // int64
	ValueTypeBool                      // bool
	ValueTypeFloat                     // float64
	ValueTypeArray                     // []any
	ValueTypeVariable                  // 变量 // string
	ValueTypeExpr                      // 表达式，函数等 // string
)

// Is 判断是否为指定类型
func (v ValueType) Is(t ValueType) bool {
	return v == t
}

// interface definition

// Parser 解析器接口
type Parser interface {

	// GetAction 获取当前query操作行为
	GetAction() Action

	// GetObject 获取第一层对象
	GetObject(key string) (object Object, exist bool)

	// Range 按顺序遍历所有 key & 对象
	Range(func(idx int, key string, object Object) (continue_ bool))

	// ToQuery 序列化
	ToQuery() (query string, err error)
}

// Object 对象接口
type Object interface {

	//
	Name() string // 对象名/虚拟字段名

	GetAlias() string

	Key() string

	Path() []string

	// get

	GetWhere() (where *Condition, exist bool)

	WalkWhere(visitor WhereVisitor) error

	GetOrderBy() (order []OrderBy, exist bool)

	IsDistinct() bool

	GetLimit() (limit int, exist bool)

	GetOffset() (offset int, exist bool)

	GetValues(dest any) error

	GetOnConflict() (oc *OnConflict, exist bool)

	GetSet(dest any) error

	GetFields() []Pair

	GetAggr() []string

	GetSub(key string) (sub Object, exist bool)

	// set

	Alias(alias string) Object

	Where(where *Condition) Object

	OrderBy(orders []OrderBy) Object

	Distinct(distinct bool) Object

	Limit(limit int) Object

	Offset(offset int) Object

	Values(dest any) Object

	OnConflict(conflict *OnConflict) Object

	Set(dest any) Object // struct | map

	Select(fields ...any) Object // string | pair | Object

	Append(fields ...any) Object // string | pair | Object

	Error() error

	// build 内部实现用
	build(parser Parser, builder *strings.Builder) error
}

// structure definition

// Condition 条件
type Condition struct {
	Connector Connector
	Exprs     []Expr
}

// Expr 表达式
type Expr struct {
	Cond  *Condition // 为nil则为普通条件
	Field string
	Op    string
	Value interface{}
	Vtype ValueType
}

// OrderBy 排序
type OrderBy struct {
	Field string
	Desc  bool
}

// OnConflict _on_conflict
type OnConflict struct {
	Columns       []string
	DoNothing     bool
	UpdateAll     bool
	UpdateColumns []string
	Update        []Pair // update { field: value }
}

// Pair { key: value } object形式的左边为K，右边为V
type Pair struct {
	K string
	V string
}

// WhereVisitor where节点遍历执行逻辑 接口
type WhereVisitor interface {
	VisitNode(key string, op string, value interface{}, kind string, logicOpStack []string) (stop bool)
}

// RawExpr 原生表达式
type RawExpr string

// Config 配置
type Config struct {
}

// function definition

// NewParserFn 创建解析器函数
type NewParserFn func(action Action, objects ...Object) (Parser, error)

// NewParserByQueryFn 创建解析器函数
type NewParserByQueryFn func(query string) (Parser, error)

// NewObjectFn 创建对象函数
type NewObjectFn func(name string) Object
