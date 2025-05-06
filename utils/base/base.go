package base

type BaseCRUDInterface interface {
	Create(data map[string]interface{}, identifier string) (bool, interface{})
	Read(query interface{}) (bool, interface{})
	Update(identifier string, newData map[string]interface{}) (bool, interface{})
	Delete(identifier string) (bool, interface{})
	Stats(query interface{}) (bool, interface{})
	Run()
}