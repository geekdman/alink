package base

type BaseCRUDInterface interface {
	Create(data interface{}, identifier string) (bool, interface{})
	Read(query interface{}) (bool, interface{})
	Update(identifier string, newData interface{}) (bool, interface{})
	Delete(identifier string) (bool, interface{})
	Stats() (bool, interface{})
	Close()
	Run()
}
// BaseCRUDInterface 定义了一个抽象接口，包含CRUD操作和其他方法
//type BaseCRUDInterface interface {
//	Create(...interface{})
//	Read(...interface{})
//	Update(...interface{})
//	Delete(...interface{})
//	Stats(...interface{})
//	Close()
//	Run()
//}