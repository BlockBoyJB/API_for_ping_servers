package database

import (
	"database/sql"
)

const (
	CreatePingData = "create table if not exists pingdata(id integer primary key AUTOINCREMENT, api_key varchar not null,  email varchar not null, url varchar not null)"
)

type (
	PingData struct {
		db     *sql.DB
		Id     int
		ApiKey string
		Email  string
		Url    string
	}
)

func initPingData(db *sql.DB) *PingData {
	return &PingData{db: db}
}

func (p *PingData) createTable() error {
	_, err := p.db.Exec(CreatePingData)
	return err
}

func (p *PingData) Create(apiKey string, email string, url string) error {
	_, err := p.db.Exec("insert into pingdata (api_key, email, url) values (?, ?, ?)", apiKey, email, url)
	return err
}

func (p *PingData) Delete(apiKey string, url string) error {
	_, err := p.db.Exec("delete from pingdata where api_key = ? and url = ?", apiKey, url)
	return err
}

func (p *PingData) Get(apiKey string, url string) (PingData, error) {
	var ping PingData
	err := p.db.QueryRow("select * from pingdata where api_key = ? and url = ?", apiKey, url).
		Scan(&ping.Id, &ping.ApiKey, &ping.Email, &ping.Url)
	return ping, err

}

func (p *PingData) GetMany(apiKey string) ([]*PingData, error) {
	data := make([]*PingData, 0)
	rows, err := p.db.Query("select * from pingdata where api_key = ?", apiKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pingdata PingData
		if err = rows.Scan(&pingdata.Id, &pingdata.ApiKey, &pingdata.Email, &pingdata.Url); err != nil {
			return nil, err
		}
		data = append(data, &pingdata)
	}
	return data, nil
}

func (p *PingData) GetAll() ([]*PingData, error) {
	data := make([]*PingData, 0)
	rows, err := p.db.Query("select * from pingdata")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var pingdata PingData
		if err = rows.Scan(&pingdata.Id, &pingdata.ApiKey, &pingdata.Email, &pingdata.Url); err != nil {
			return nil, err
		}
		data = append(data, &pingdata)
	}
	return data, nil
}
