package skud

const (
	AccessGranted                  = "доступ дозволений"
	AccessDeniedInaccessible       = "доступ заборонений: користувач поза зоною зчитувача"
	AccessGrantedWithHealthCheck   = "доступ дозволений з відмітккою у лікаря"
	AccessDeniedNoHealthCheck      = "доступ заборонений: необхідна відмітка у лікаря"
	AccessGrantedWithSanitaryCheck = "доступ дозволений, актуальний санітарний одяг"
	AccessDeniedNoSanitaryCheck    = "доступ заборонений: застарілий санітарний одяг"
	AccessGrantedWithAllChecks     = "доступ дозволений з відміткою у лікаря та актуальним санітарним одягом"
	AccessDeniedNoAnyChecks        = "доступ заборонений: необхідна відмітка у лікаря та актуальний санітарний одяг"

	AccessDeniedUnknownEmployee = "доступ заборонений: картку не розпізнано"
	AccessDeniedWrongPath       = "доступ заборонений: невірний шлях"
)
