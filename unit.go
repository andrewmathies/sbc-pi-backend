package bg

type Unit struct {
	Version string 	`json:"version"`
	BeanID 	string 	`json:"beanID"`
	Name	string	`json:"name"`
	State	State	`json:"state"`
}