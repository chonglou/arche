package nut

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// H hash
type H map[string]interface{}

// Locale locale
type Locale struct {
	ID        uint      `json:"id" orm:"column(id)"`
	Lang      string    `json:"lang"`
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	UpdatedAt time.Time `json:"updatedAt" orm:"auto_now"`
	CreatedAt time.Time `json:"createdAt" orm:"auto_now_add"`
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
	UpdatedAt time.Time `json:"updatedAt" orm:"auto_now"`
	CreatedAt time.Time `json:"createdAt" orm:"auto_now_add"`
}

// TableName table name
func (p *Setting) TableName() string {
	return "settings"
}

// -----------------------------------------------------------------------------

// User user
type User struct {
	ID              uint       `json:"id" orm:"column(id)"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	UID             string     `json:"uid" orm:"column(uid)"`
	Password        string     `json:"-"`
	ProviderID      string     `json:"providerId" orm:"column(provider_id)"`
	ProviderType    string     `json:"providerType"`
	Logo            string     `json:"logo"`
	SignInCount     uint       `json:"signInCount"`
	LastSignInAt    *time.Time `json:"lastSignInAt"`
	LastSignInIP    string     `json:"lastSignInIp" orm:"column(last_sign_in_ip)"`
	CurrentSignInAt *time.Time `json:"currentSignInAt"`
	CurrentSignInIP string     `json:"currentSignInIp" orm:"column(current_sign_in_ip)"`
	ConfirmedAt     *time.Time `json:"confirmedAt"`
	LockedAt        *time.Time `json:"lockedAt"`
	UpdatedAt       time.Time  `json:"updatedAt" orm:"auto_now"`
	CreatedAt       time.Time  `json:"createdAt" orm:"auto_now_add"`

	Logs []*Log `json:"logs" orm:"reverse(many)"`
}

// SetEmail set email
func (p *User) SetEmail(s string) {
	p.Email = strings.ToLower(s)
}

// SetPassword set password
func (p *User) SetPassword(s string) error {
	buf, err := bcrypt.GenerateFromPassword([]byte(s), 16)
	if err != nil {
		return err
	}
	p.Password = string(buf)
	return nil
}

// Auth check email & password
func (p *User) Auth(email, password string) bool {
	return p.ProviderType == UserTypeEmail &&
		strings.ToLower(email) == p.Email &&
		bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(password)) == nil
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

func (p User) String() string {
	return fmt.Sprintf("%s<%s>", p.Name, p.Email)
}

// TableName table name
func (p *User) TableName() string {
	return "users"
}

// Attachment attachment
type Attachment struct {
	ID           uint      `json:"id" orm:"column(id)"`
	Title        string    `json:"title"`
	URL          string    `json:"url" orm:"column(url)"`
	Length       int64     `json:"length"`
	MediaType    string    `json:"mediaType"`
	ResourceID   uint      `json:"resourceId" orm:"column(resource_id)"`
	ResourceType string    `json:"resourceType"`
	User         *User     `json:"user" orm:"rel(fk)"`
	UpdatedAt    time.Time `json:"updatedAt" orm:"auto_now"`
	CreatedAt    time.Time `json:"createdAt" orm:"auto_now_add"`
}

// IsPicture is picture?
func (p *Attachment) IsPicture() bool {
	return strings.HasPrefix(p.MediaType, "image/")
}

// TableName table name
func (p *Attachment) TableName() string {
	return "attachments"
}

// Log log
type Log struct {
	ID        uint      `json:"id" orm:"column(id)"`
	Message   string    `json:"message"`
	IP        string    `json:"ip" orm:"column(ip)"`
	User      *User     `json:"user" orm:"rel(fk)"`
	CreatedAt time.Time `json:"createdAt" orm:"auto_now_add"`
}

func (p Log) String() string {
	return fmt.Sprintf("%s: [%s]\t %s", p.CreatedAt.Format(time.ANSIC), p.IP, p.Message)
}

// TableName table name
func (p *Log) TableName() string {
	return "logs"
}

// Policy policy
type Policy struct {
	ID        uint `orm:"column(id)"`
	Nbf       time.Time
	Exp       time.Time
	User      *User     `orm:"rel(fk)"`
	Role      *Role     `orm:"rel(fk)"`
	UpdatedAt time.Time `json:"updatedAt" orm:"auto_now"`
	CreatedAt time.Time `json:"createdAt" orm:"auto_now_add"`
}

//Enable is enable?
func (p *Policy) Enable() bool {
	now := time.Now()
	return now.After(p.Nbf) && now.Before(p.Exp)
}

// TableName table name
func (p *Policy) TableName() string {
	return "policies"
}

// Role role
type Role struct {
	ID           uint `orm:"column(id)"`
	Name         string
	ResourceID   uint `orm:"column(resource_id)"`
	ResourceType string
	UpdatedAt    time.Time `json:"updatedAt" orm:"auto_now"`
	CreatedAt    time.Time `json:"createdAt" orm:"auto_now_add"`
}

func (p Role) String() string {
	return fmt.Sprintf("%s@%s://%d", p.Name, p.ResourceType, p.ResourceID)
}

// TableName table name
func (p *Role) TableName() string {
	return "roles"
}

// -----------------------------------------------------------------------------

// Vote vote
type Vote struct {
	ID           uint `orm:"column(id)"`
	Point        int
	ResourceID   uint `orm:"column(resource_id)"`
	ResourceType string
	UpdatedAt    time.Time `json:"updatedAt" orm:"auto_now"`
	CreatedAt    time.Time `json:"createdAt" orm:"auto_now_add"`
}

// TableName table name
func (p Vote) TableName() string {
	return "votes"
}

// LeaveWord leave-word
type LeaveWord struct {
	ID        uint      `json:"id" orm:"column(id)"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt" orm:"auto_now_add"`
}

// TableName table name
func (p LeaveWord) TableName() string {
	return "leave_words"
}

// Link link
type Link struct {
	ID        uint      `json:"id" orm:"column(id)"`
	Lang      string    `json:"lang"`
	X         int       `json:"x"`
	Y         int       `json:"y"`
	Href      string    `json:"href"`
	Label     string    `json:"label"`
	Loc       string    `json:"loc"`
	UpdatedAt time.Time `json:"updatedAt" orm:"auto_now"`
	CreatedAt time.Time `json:"createdAt" orm:"auto_now_add"`
}

// TableName table name
func (p *Link) TableName() string {
	return "links"
}

// Card card
type Card struct {
	ID        uint      `json:"id" orm:"column(id)"`
	Lang      string    `json:"lang"`
	Loc       string    `json:"loc"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Type      string    `json:"type"`
	Href      string    `json:"href"`
	Logo      string    `json:"logo"`
	Sort      int       `json:"sort"`
	Action    string    `json:"action"`
	UpdatedAt time.Time `json:"updatedAt" orm:"auto_now"`
	CreatedAt time.Time `json:"createdAt" orm:"auto_now_add"`
}

// TableName table name
func (p *Card) TableName() string {
	return "cards"
}

// FriendLink friend_links
type FriendLink struct {
	ID        uint      `json:"id" orm:"column(id)"`
	Title     string    `json:"title"`
	Home      string    `json:"home"`
	Logo      string    `json:"logo"`
	Sort      int       `json:"sort"`
	UpdatedAt time.Time `json:"updatedAt" orm:"auto_now"`
	CreatedAt time.Time `json:"createdAt" orm:"auto_now_add"`
}

// TableName table name
func (p *FriendLink) TableName() string {
	return "friend_links"
}

// -----------------------------------------------------------------------------

func init() {
	orm.RegisterModel(
		new(Locale), new(Setting),
		new(User), new(Log), new(Role), new(Policy),
		new(Card), new(Link),
		new(LeaveWord), new(Attachment),
		new(FriendLink), new(Vote),
	)
}
