package entities

import (
	"io"
	"mime/multipart"
	"os"
	"supinfo/mewpipe/configs"

	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func getMediaCollection() *mgo.Collection {
	return configs.MongoDB.C("users")
}

func getMediaGridFSCollection() *mgo.GridFS {
	return configs.MongoDB.GridFS("media")
}

type scope string

const (
	Public  scope = "public"  //Available to anybody
	Private scope = "private" //Available to authenticated users
	Link    scope = "link"    //Available to anybody with the link
)

type user struct {
	Id    bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name  name          `json:"name" bson:",omitempty"`
	Email string        `json:"email" bson:",omitempty"`
}

type Media struct {
	Id        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	Title     string        `json:"title" bson:",omitempty"`
	Summary   string        `json:"summary" bson:",omitempty"`
	Publisher user          `json:"user,omitempty" bson:",omitempty"`
	File      bson.ObjectId `json:"file,omitempty" bson:",omitempty"`
	Scope     scope         `json:"scope,omitempty" bson:"scope,omitempty"`
	mgofile   *mgo.GridFile `json:"-" bson:"-"`
}

func MediaNew() *Media {
	media := new(Media)
	media.Id = bson.NewObjectId()
	media.CreatedAt = time.Now()
	return media
}

func MediaNewFromId(oid bson.ObjectId) *Media {
	media := MediaNew()
	media.Id = oid
	return media
}
func MediaFromId(oid bson.ObjectId) (*Media, error) {
	media := new(Media)
	err := getMediaCollection().FindId(oid).One(&media)
	if err != nil {
		media = MediaNewFromId(oid)
	}
	return media, err
}

func MediaList(bson bson.M, start int, number int, sort ...string) ([]Media, error) {
	medias := make([]Media, number)

	err := getMediaCollection().Find(bson).Sort(sort...).Skip(start).Limit(number).All(&medias)

	return medias, err
}

func (m *Media) Normalize() {
	if m.Scope != "link" && m.Scope != "private" {
		m.Scope = "public"
	}
}

func (m *Media) Upload(postedFile io.Reader, fileHeader *multipart.FileHeader) error {
	mongoFile, err := getMediaGridFSCollection().Create(fileHeader.Filename)
	defer mongoFile.Close()
	if err != nil {
		return err
	}
	mongoFile.SetContentType(fileHeader.Header.Get("Content-Type"))
	_, err = io.Copy(mongoFile, postedFile)
	if err != nil {
		return err
	}

	if m.File != "" {
		getMediaGridFSCollection().RemoveId(m.File)
	}
	m.File = mongoFile.Id().(bson.ObjectId)

	if err := m.Update(); err != nil {
		mongoFile.Abort()
		return err
	}
	return nil
}

func (m *Media) OpenFile() error {
	file, err := getMediaGridFSCollection().OpenId(m.File)
	m.mgofile = file

	if err != nil {
		return err
	}

	return nil
}

func (m *Media) ContentType() string {
	return m.mgofile.ContentType()
}

func (m *Media) Size() int64 {
	return m.mgofile.Size()
}

func (m *Media) CloseFile() error {
	return m.mgofile.Close()
}

func (m *Media) SeekSet(offset int64) error {
	_, err := m.mgofile.Seek(offset, os.SEEK_SET)
	return err
}

func (m *Media) Read(buffer []byte) error {
	_, err := m.mgofile.Read(buffer)
	return err
}

func (m *Media) CopyTo(target io.Writer) error {
	_, err := io.Copy(target, m.mgofile)
	return err
}

func (m *Media) Insert() error {
	m.Normalize()
	if err := getMediaCollection().Insert(&m); err != nil {
		return err
	}
	return nil
}

func (m *Media) Update() error {
	m.Normalize()
	if err := getMediaCollection().UpdateId(m.Id, &m); err != nil {
		return err
	}
	return nil
}

func (m *Media) Delete() error {
	if m.File != "" {
		getMediaGridFSCollection().RemoveId(m.File)
	}
	if err := getMediaCollection().RemoveId(m.Id); err != nil {
		return err
	}
	return nil
}
