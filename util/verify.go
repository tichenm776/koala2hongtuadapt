package util

var (
	//IdVerify               = Rules{"ID": {NotEmpty()}}
	//ApiVerify              = Rules{"Path": {NotEmpty()}, "Description": {NotEmpty()}, "ApiGroup": {NotEmpty()}, "Method": {NotEmpty()}}
	//MenuVerify             = Rules{"Path": {NotEmpty()}, "ParentId": {NotEmpty()}, "Name": {NotEmpty()}, "Component": {NotEmpty()}, "Sort": {Ge("0")}}
	//MenuMetaVerify         = Rules{"Title": {NotEmpty()}}
	//LoginVerify            = Rules{"username": {NotEmpty()}, "password": {NotEmpty()}}
	//FileVerify            = Rules{"file_path": {NotEmpty()}, "file_name": {NotEmpty()}, "file_number": {NotEmpty()}, "caseid": {NotEmpty()}, "userid": {NotEmpty()}, "filesid": {NotEmpty()}}
	//RegisterVerify         = Rules{"username": {NotEmpty()}, "phone": {NotEmpty()}, "password": {NotEmpty()}}
	//PageInfoVerify         = Rules{"Page": {NotEmpty()}, "PageSize": {NotEmpty()}}
	//AddCaseVerify         = Rules{"invoice_handover": {NotEmpty()},"note": {NotEmpty()},"receive_department": {NotEmpty()},"litigant": {NotEmpty()},"userid": {NotEmpty()}, "lawyername": {NotEmpty()}, "status": {NotEmpty()}, "money": {NotEmpty()}, "createtime": {NotEmpty()}, "casetype": {NotEmpty()}, "casebrief": {NotEmpty()}}
	//DeleteCaseVerify         = Rules{"userid": {NotEmpty()}, "caseid": {NotEmpty()}}
	//EdificesDTOVerify         = Rules{"EdificeName": {NotEmpty()}}
	//EdificesupdateDTOVerify         = Rules{"EdificeName": {NotEmpty()},"EdificesId":{NotEmpty()}}
	//EdificesdeleteDTOVerify         = Rules{"EdificesId":{NotEmpty()}}
	//FloorsupdateDTOVerify         = Rules{"FloorId": {NotEmpty()},"FloorName":{NotEmpty()},"EdificeName":{NotEmpty()},"EdificeId":{NotEmpty()}}
	//FloorsdeleteDTOVerify         = Rules{"FloorId": {NotEmpty()},"EdificeId":{NotEmpty()}}
	//FloorsDTOVerify         = Rules{"FloorName":{NotEmpty()},"EdificeName":{NotEmpty()},"EdificeId":{NotEmpty()}}
	//CompanyDTOVerify         = Rules{"Name":{NotEmpty()},"Locations":{NotEmpty()}}
	//CompanyDTOupdateVerify         = Rules{"ID":{NotEmpty()},"Name":{NotEmpty()},"Locations":{NotEmpty()}}
	//CompanyDTOdeleteVerify         = Rules{"ID":{NotEmpty()}}
	//
	//DepartmentCreateDTOVerify         = Rules{"Name":{NotEmpty()},"CompanyId":{NotEmpty()},"ParentId":{NotEmpty()}}
	//DepartmentUpdateDTOVerify         = Rules{"ID":{NotEmpty()},"Name":{NotEmpty()},"CompanyId":{NotEmpty()},"ParentId":{NotEmpty()}}
	////DepartmentDTOupdateVerify         = Rules{"ID":{NotEmpty()},"Name":{NotEmpty()},"Locations":{NotEmpty()}}
	//DepartmentDTOdeleteVerify         = Rules{"ID":{NotEmpty()}}
	//
	//EmployeeCreateDTOVerify	=Rules{"Name":{NotEmpty()},"CompanyId":{NotEmpty()},"FloorId":{NotEmpty()},"Phone":{NotEmpty()}}
	//EmployeeUpdateDTOVerify	=Rules{"Name":{NotEmpty()},"CompanyId":{NotEmpty()},"ID":{NotEmpty()},"Phone":{NotEmpty()}}
	//EmployeeListDTOVerify	=Rules{"CompanyId":{NotEmpty()}}
	//
	//VisitorCreateDTOVerify	=Rules{"Name":{NotEmpty()},
	//	"Gender":{NotEmpty()},"Phone":{NotEmpty()},
	//	"VisitTime":{NotEmpty()},"VisitReason":{NotEmpty()},
	//	"AuthTimes":{NotEmpty()},"StartTime":{NotEmpty()},
	//	"EndTime":{NotEmpty()},"VisitType":{NotEmpty()}}


	SubjectCreateDTOVerify	=Rules{"Subject_type":{NotEmpty()},
		"Name":{NotEmpty()}}





	//LocationsCompanyDTOVerify         = Rules{"EdificeId":{NotEmpty()},"FloorId":{NotEmpty()}}

	//CustomerVerify         = Rules{"CustomerName": {NotEmpty()}, "CustomerPhoneData": {NotEmpty()}}
	//AutoCodeVerify         = Rules{"Abbreviation": {NotEmpty()}, "StructName": {NotEmpty()}, "PackageName": {NotEmpty()}, "Fields": {NotEmpty()}}
	//AuthorityVerify        = Rules{"AuthorityId": {NotEmpty()}, "AuthorityName": {NotEmpty()}, "ParentId": {NotEmpty()}}
	//AuthorityIdVerify      = Rules{"AuthorityId": {NotEmpty()}}
	//OldAuthorityVerify     = Rules{"OldAuthorityId": {NotEmpty()}}
	//ChangePasswordVerify   = Rules{"Username": {NotEmpty()}, "Password": {NotEmpty()}, "NewPassword": {NotEmpty()}}
	//SetUserAuthorityVerify = Rules{"UUID": {NotEmpty()}, "AuthorityId": {NotEmpty()}}
)