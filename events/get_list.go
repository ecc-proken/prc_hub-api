package events

import (
	"prc_hub-api/mysql"
	"prc_hub-api/users"
	"strings"
	"time"
)

type GetQuery struct {
	Published       *bool   `json:"published" validate:"omitempty"`
	Title           *string `json:"name" validate:"omitempty"`
	TitleContain    *string `json:"name_contain" validate:"omitempty"`
	Location        *string `json:"location" validate:"omitempty"`
	LocationContain *string `json:"location_contain" validate:"omitempty"`
}

func Get(query GetQuery, userId *uint64, admin bool) (events []Event, err error) {
	// クエリを作成
	queryStrBase :=
		`SELECT
			id AS event_id
		FROM events
		WHERE`
	queryParams := []interface{}{}

	if userId != nil && !admin {
		queryStrBase += " published = true OR user_id = ? AND"
		queryParams = append(queryParams, userId)
	}
	if userId == nil {
		queryStrBase += " published = true AND"
	}
	if userId != nil && query.Published != nil {
		queryStrBase += " published = ? AND"
		queryParams = append(queryParams, query.Published)
	}
	if query.Title != nil {
		queryStrBase += " title = ? AND"
		queryParams = append(queryParams, query.Title)
	}
	if query.TitleContain != nil {
		queryStrBase += " title LIKE ? AND"
		queryParams = append(queryParams, "%"+*query.TitleContain+"%")
	}
	if query.Location != nil {
		queryStrBase += " location = ? AND"
		queryParams = append(queryParams, query.Location)
	}
	if query.LocationContain != nil {
		queryStrBase += " location LIKE ?"
		queryParams = append(queryParams, "%"+*query.LocationContain+"%")
	}

	queryStrBase = strings.TrimSuffix(queryStrBase, "WHERE")
	queryStrBase = strings.TrimSuffix(queryStrBase, "AND")

	queryStr :=
		`WITH params AS (
			` + queryStrBase + `
		)
		SELECT
			e.id, e.title, e.description, e.location, e.published, e.completed,
			null, null, null, null,
			null, null, null, null,
			null, null, null
		FROM events e
		WHERE e.id IN (SELECT event_id FROM params)
		UNION ALL
		SELECT
			s.event_id, null, null, null, null, null,
			u.id, u.name, u.github_username, u.twitter_id,
			null, null, null, null,
			null, null, null
		FROM event_speakers s, users u
		WHERE s.event_id IN (SELECT event_id FROM params) AND s.user_id = u.id
		UNION ALL
		SELECT
			dt.event_id, null, null, null, null, null,
			null, null, null, null,
			dt.id, dt.start, dt.end, ep.count,
			null, null, null
		FROM event_datetimes dt
		LEFT JOIN
			(SELECT event_datetime_id, COUNT(event_datetime_id) AS count FROM event_participates GROUP BY event_datetime_id) AS ep
		ON dt.id = ep.event_datetime_id
		WHERE dt.event_id IN (SELECT event_id FROM params)
		UNION ALL
		SELECT
			doc.event_id, null, null, null, null, null,
			null, null, null, null,
			null, null, null, null,
			doc.id, doc.name, doc.url
		FROM event_documents doc
		WHERE doc.event_id IN (SELECT event_id FROM params)
		ORDER BY id`

	rows, err := mysql.Read(queryStr, queryParams...)
	if err != nil {
		return
	}
	defer rows.Close()

	// 読込中Event
	var loadingEvent *Event
	// 1行ずつ読込
	for rows.Next() {
		// 読込用変数
		var (
			eId          uint64
			eTitle       *string
			eDescription *string
			eLocation    *string
			ePublished   *bool
			eCompleted   *bool

			uId      *uint64
			uName    *string
			uGithub  *string
			uTwitter *string

			dtId    *uint64
			dtStart *time.Time
			dtEnd   *time.Time
			dpCount *uint

			dcId   *uint64
			dcName *string
			dcUrl  *string
		)
		// 変数に割り当て
		err = rows.Scan(
			&eId, &eTitle, &eDescription, &eLocation, &ePublished, &eCompleted,
			&uId, &uName, &uGithub, &uTwitter,
			&dtId, &dtStart, &dtEnd, &dpCount,
			&dcId, &dcName, &dcUrl,
		)

		if loadingEvent == nil {
			// 読込中のEventがない場合(初回に実行)
			// 新しく読み込んだEventを保持
			loadingEvent = &Event{
				Id:          eId,
				Title:       *eTitle,
				Description: eDescription,
				Location:    eLocation,
				Published:   *ePublished,
				Completed:   *eCompleted,
			}
		} else if eId != loadingEvent.Id {
			// Eventが変わった場合
			// レスポンス用の配列に追加
			events = append(events, *loadingEvent)

			// 新しく読み込んだEventを保持
			loadingEvent = &Event{
				Id:          eId,
				Title:       *eTitle,
				Description: eDescription,
				Location:    eLocation,
				Published:   *ePublished,
				Completed:   *eCompleted,
			}
		} else if uId != nil && uName != nil {
			// UserをEvent.Speakersに追加
			loadingEvent.Speakers = append(
				loadingEvent.Speakers,
				users.UserEmbed{
					Id:             *uId,
					Name:           *uName,
					GithubUsername: uGithub,
					TwitterId:      uTwitter,
				},
			)
		} else if dtId != nil && dtStart != nil {
			// EventDatetimeをEvent.Datetimesに追加
			loadingEvent.Datetimes = append(
				loadingEvent.Datetimes,
				EventDatetime{
					Id:               *dtId,
					EventId:          eId,
					Start:            *dtStart,
					End:              dtEnd,
					ParticipateCount: dpCount,
				},
			)
		} else if dcId != nil && dcName != nil && dcUrl != nil {
			// EventDocumentをEvent.Documentsに追加
			loadingEvent.Documents = append(
				loadingEvent.Documents,
				EventDocument{
					Id:   *dcId,
					Name: *dcName,
					Url:  *dcUrl,
				},
			)
		}
	}
	if loadingEvent != nil {
		events = append(events, *loadingEvent)
	}

	return
}
