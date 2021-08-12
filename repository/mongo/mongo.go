package mongo
import (
	"context"
	"sync"

	"awesomeProject6/model"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type mongoRepository struct {
	session *mgo.Session
}

var defaultDBName = "API"
var defaultCollection = "students"
var userCollection = "user"

var mutex sync.Mutex
var mongoRepo = &mongoRepository{}

func Init(host string) *mongoRepository {
	mutex.Lock()
	defer mutex.Unlock()

	if mongoRepo.session != nil {
		return mongoRepo
	}

	session, err := mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	mongoRepo.session = session
	return mongoRepo
}

func (repo *mongoRepository) CreateStudent(ctx context.Context, Student *model.StudentDetails) (*model.StudentDetails, error) {
	session := repo.session.Copy()
	defer session.Close()

	if Student.Id == "" {
		Student.Id = bson.NewObjectId()
	}

	err := session.DB(defaultDBName).C(defaultCollection).Insert(Student)
	if err != nil {
		return nil, err
	}
	return Student, nil
}

func (repo *mongoRepository) CreateUser(ctx context.Context, User *model.Credentials) (*model.Credentials, error) {
	session := repo.session.Copy()
	defer session.Close()

	if User.Id == "" {
		User.Id = bson.NewObjectId()
	}

	err := session.DB(defaultDBName).C(userCollection).Insert(User)
	if err != nil {
		return nil, err
	}
	return User, nil
}

func (repo *mongoRepository) DeleteStudent(ctx context.Context, id string) error {
	session := repo.session.Copy()
	defer session.Close()

	m := bson.M{"_id": bson.ObjectIdHex(id)}

	return session.DB(defaultDBName).C(defaultCollection).Remove(m)
}

func (repo *mongoRepository) UpdateStudent(ctx context.Context,
	Student *model.StudentDetails) (*model.StudentDetails, error) {
	session := repo.session.Copy()
	defer session.Close()

	m := bson.M{"_id": Student.Id}

	err := session.DB(defaultDBName).C(defaultCollection).Update(m, Student)
	if err != nil {
		return nil, err
	}

	return Student, nil
}

func (repo *mongoRepository) GetStudent(ctx context.Context, id string) (*model.StudentDetails,
	error) {
	session := repo.session.Copy()
	defer session.Close()

	m := bson.M{"_id": bson.ObjectIdHex(id)}

	var Student = &model.StudentDetails{}
	collection := session.DB(defaultDBName).C(defaultCollection)
	err := collection.Find(m).One(Student)
	if err != nil {
		return nil, err
	}
	return Student, nil
}

func (repo *mongoRepository) GetUser(ctx context.Context, id string) (*model.Credentials,
	error) {
	session := repo.session.Copy()
	defer session.Close()

	m := bson.M{"_id": bson.ObjectIdHex(id)}


	var User = &model.Credentials{}
	collection := session.DB(defaultDBName).C(userCollection)
	err := collection.Find(m).One(User)
	if err != nil {
		return nil, err
	}
	return User, nil
}

func (repo *mongoRepository) ListStudent(ctx context.Context) ([]*model.StudentDetails, error) {

	session := repo.session.Copy()
	defer session.Close()
	var Students []*model.StudentDetails

	m := bson.M{}
	err := session.DB(defaultDBName).C(defaultCollection).Find(m).All(&Students)

	if err != nil {
		return nil, err
	}

	return Students, nil
}

func (repo *mongoRepository) Close() {
	repo.session.Close()
}

