// Copyright (C) 2024 Declan Teevan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package user

import (
	"github.com/hexolan/stocklet/internal/pkg/messaging"
	eventspb "github.com/hexolan/stocklet/internal/pkg/protogen/events/v1"
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/user/v1"
)

func PrepareUserCreatedEvent(user *pb.User) ([]byte, string, error) {
	topic := messaging.User_State_Created_Topic
	event := &eventspb.UserCreatedEvent{
		Revision: 1,

		UserId:    user.Id,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareUserEmailUpdatedEvent(userId string, email string) ([]byte, string, error) {
	topic := messaging.User_Attribute_Email_Topic
	event := &eventspb.UserEmailUpdatedEvent{
		Revision: 1,

		UserId: userId,
		Email:  email,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareUserDeletedEvent(user *pb.User) ([]byte, string, error) {
	topic := messaging.User_State_Deleted_Topic
	event := &eventspb.UserDeletedEvent{
		Revision: 1,

		UserId: user.Id,
		Email:  user.Email,
	}

	return messaging.MarshalEvent(event, topic)
}
