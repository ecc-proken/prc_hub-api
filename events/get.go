package events

import (
	"prc_hub-api/mysql"
	"prc_hub-api/users"
	"time"
)

func GetById(id uint64) (e Event, notFound bool, err error) {
	rows, err := mysql.Read(
		`WITH params AS (
			SELECT ? as event_id
		)
		SELECT
			e.id, e.title, e.description, e.location, e.published, e.completed, e.auto_notify_documents_enabled,
			null, null, null, null,
			null, null, null,
			null, null, null
		FROM events e
		WHERE e.id IN (SELECT event_id FROM params)
		UNION ALL
		SELECT
			s.event_id, null, null, null, null, null, null,
			u.id, u.name, u.github_username, u.twitter_id,
			null, null, null,
			null, null, null
		FROM event_speakers s, users u
		WHERE s.event_id IN (SELECT event_id FROM params) AND s.user_id = u.id
		UNION ALL
		SELECT
			dt.event_id, null, null, null, null, null, null,
			null, null, null, null,
			dt.id, dt.start, dt.end,
			null, null, null
		FROM event_datetimes dt
		WHERE dt.event_id IN (SELECT event_id FROM params)
		UNION ALL
		SELECT
			doc.event_id, null, null, null, null, null, null,
			null, null, null, null,
			null, null, null,
			doc.id, doc.name, doc.url
		FROM event_documents doc
		WHERE doc.event_id IN (SELECT event_id FROM params)
		ORDER BY id`,
		id,
	)
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
			eId                  uint64
			eTitle               *string
			eDescription         *string
			eLocation            *string
			ePublished           *bool
			eCompleted           *bool
			eAutoNotifyDocuments *bool

			uId      *uint64
			uName    *string
			uGithub  *string
			uTwitter *string

			dtId    *uint64
			dtStart *time.Time
			dtEnd   *time.Time

			dcId   *uint64
			dcName *string
			dcUrl  *string
		)
		// 変数に割り当て
		err = rows.Scan(
			&eId, &eTitle, &eDescription, &eLocation, &ePublished, &eCompleted, &eAutoNotifyDocuments,
			&uId, &uName, &uGithub, &uTwitter,
			&dtId, &dtStart, &dtEnd,
			&dcId, &dcName, &dcUrl,
		)

		if loadingEvent == nil {
			// 読込中のEventがない場合(初回に実行)
			// 新しく読み込んだEventを保持
			loadingEvent = &Event{
				Id:                  eId,
				Title:               *eTitle,
				Description:         eDescription,
				Location:            *eLocation,
				Published:           *ePublished,
				Completed:           *eCompleted,
				AutoNotifyDocuments: *eAutoNotifyDocuments,
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
					Id:      *dtId,
					EventId: eId,
					Start:   *dtStart,
					End:     dtEnd,
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

	if loadingEvent == nil {
		notFound = true
		return
	} else {
		e = *loadingEvent
	}
	return
}
