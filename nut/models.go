package nut

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/google/uuid"
)

// Locale locale
type Locale struct {
	ID        uint      `json:"id" orm:"column(id)"`
	Lang      string    `json:"lang"`
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// TableName table name
func (p *Locale) TableName() string {
	return "locales"
}

// -----------------------------------------------------------------------------

// Setting setting model
type Setting struct {
	ID        uint `orm:"column(id)"`
	Key       string
	Value     string
	Encode    bool
	UpdatedAt time.Time
	CreatedAt time.Time
}

// TableName table name
func (p *Setting) TableName() string {
	return "settings"
}

// -----------------------------------------------------------------------------

// Domain domain
type Domain struct {
	ID        uint      `json:"id" orm:"column(id)"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

//SetName set name
func (p *Domain) SetName(n string) {
	p.Name = strings.ToLower(n)
}

// TableName table name
func (p *Domain) TableName() string {
	return "domains"
}

// User user
type User struct {
	ID              uint       `json:"id" orm:"column(id)"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	UID             string     `json:"uid" orm:"column(uid)"`
	Password        string     `json:"-"`
	Logo            string     `json:"logo"`
	SignInCount     uint       `json:"signInCount"`
	LastSignInAt    *time.Time `json:"lastSignInAt"`
	LastSignInIP    string     `json:"lastSignInIp" orm:"column(last_sign_in_ip)"`
	CurrentSignInAt *time.Time `json:"currentSignInAt"`
	CurrentSignInIP string     `json:"currentSignInIp" orm:"column(current_sign_in_ip)"`
	ConfirmedAt     *time.Time `json:"confirmedAt"`
	LockedAt        *time.Time `json:"lockAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	CreatedAt       time.Time  `json:"createdAt"`

	Logs   []*Log  `orm:"reverse(many)"`
	Domain *Domain `orm:"rel(fk)"`
}

// IsConfirm is confirm?
func (p *User) IsConfirm() bool {
	return p.ConfirmedAt != nil
}

// IsLock is lock?
func (p *User) IsLock() bool {
	return p.LockedAt != nil
}

// SetGravatarLogo set logo by gravatar
func (p *User) SetGravatarLogo() {
	// https: //en.gravatar.com/site/implement/
	buf := md5.Sum([]byte(strings.ToLower(p.Email)))
	p.Logo = fmt.Sprintf("https://gravatar.com/avatar/%s.png", hex.EncodeToString(buf[:]))
}

//SetUID generate uid
func (p *User) SetUID() {
	p.UID = uuid.New().String()
}

//SetEmail set email
func (p *User) SetEmail(e string) {
	p.Email = strings.ToLower(e)
}

func (p User) String() string {
	return fmt.Sprintf("%s<%s>", p.Name, p.Email)
}

// TableName table name
func (p *User) TableName() string {
	return "users"
}

// Alias alias
type Alias struct {
	ID          uint      `json:"id" orm:"column(id)"`
	Source      string    `json:"source"`
	Destination string    `json:"destination"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`

	Domain *Domain `orm:"rel(fk)"`
}

// SetSource set source
func (p *Alias) SetSource(s string) {
	p.Source = strings.ToLower(s)
}

// SetDestination set destination
func (p *Alias) SetDestination(d string) {
	p.Destination = strings.ToLower(d)
}

// TableName table name
func (p *Alias) TableName() string {
	return "aliases"
}

// Log log
type Log struct {
	ID        uint      `json:"id" orm:"column(id)"`
	Message   string    `json:"message"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"createdAt"`

	User *User `orm:"rel(fk)"`
}

func (p *Log) String() string {
	return fmt.Sprintf("%s: [%s]\t %s", p.CreatedAt.Format(time.ANSIC), p.IP, p.Message)
}

// TableName table name
func (p *Log) TableName() string {
	return "logs"
}

// Policy policy
type Policy struct {
	ID        uint ` orm:"column(id)"`
	StartUp   time.Time
	ShutDown  time.Time
	User      *User `orm:"reverse(one)"`
	Role      *Role `orm:"reverse(one)"`
	UpdatedAt time.Time
	CreatedAt time.Time
}

//Enable is enable?
func (p *Policy) Enable() bool {
	now := time.Now()
	return now.After(p.StartUp) && now.Before(p.ShutDown)
}

// TableName table name
func (p *Policy) TableName() string {
	return "policies"
}

// Role role
type Role struct {
	ID           uint `orm:"column(id)"`
	Name         string
	ResourceID   uint
	ResourceType string
	UpdatedAt    time.Time
	CreatedAt    time.Time
}

func (p *Role) String() string {
	return fmt.Sprintf("%s@%s://%d", p.Name, p.ResourceType, p.ResourceID)
}

// TableName table name
func (p *Role) TableName() string {
	return "roles"
}

// -----------------------------------------------------------------------------

func init() {
	orm.RegisterModel(
		new(Locale), new(Setting),
		new(Domain), new(User), new(Alias), new(Log), new(Role), new(Policy),
	)
}
