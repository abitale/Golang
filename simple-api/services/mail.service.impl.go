package services

import (
	"context"
	"errors"

	"example.com/simple-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MailServiceImpl struct {
	mailCollection *mongo.Collection
	ctx            context.Context //bisa digunakan untuk mengatur waktu yang dibutuhkan untuk melakukan proses
}

func NewMailService(mailCollection *mongo.Collection, ctx context.Context) MailService {
	return &MailServiceImpl{
		mailCollection: mailCollection,
		ctx:            ctx,
	}
}

func (m *MailServiceImpl) CreateMail(mail *models.Mail) error {
	_, err := m.mailCollection.InsertOne(m.ctx, mail)
	return err
}

func (m *MailServiceImpl) GetMail(id *int) (*models.Mail, error) {
	var mail *models.Mail
	query := bson.D{bson.E{Key: "id", Value: id}}
	err := m.mailCollection.FindOne(m.ctx, query).Decode(&mail)
	return mail, err
}

func (m *MailServiceImpl) GetAll() ([]*models.Mail, error) {
	var mails []*models.Mail
	cursor, err := m.mailCollection.Find(m.ctx, bson.D{{}})

	if err != nil {
		return nil, err
	}

	for cursor.Next(m.ctx) {
		var mail models.Mail
		err := cursor.Decode(&mail)
		if err != nil {
			return nil, err
		}
		mails = append(mails, &mail)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(m.ctx)

	if len(mails) == 0 {
		return nil, errors.New("no document found")
	}

	return mails, nil
}

func (m *MailServiceImpl) UpdateMail(id *int, mail *models.Mail) error {
	filter := bson.D{bson.E{Key: "id", Value: id}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "sender", Value: mail.Sender}, bson.E{Key: "receiver", Value: mail.Receiver}, bson.E{Key: "body", Value: mail.Body}}}}
	result, _ := m.mailCollection.UpdateOne(m.ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("no matched data found")
	}
	return nil
}

func (m *MailServiceImpl) DeleteMail(id *int) error {
	filter := bson.D{bson.E{Key: "id", Value: id}}
	result, _ := m.mailCollection.DeleteOne(m.ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("no matched data found")
	}
	return nil
}
