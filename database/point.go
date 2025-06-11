package database

import "fmt"

func (dg *DatabaseGame) addPoint() (*Point, error) {
	var point *Point

	err := dg.db.Create(point).Error

	if err != nil {
		return nil, err
	}

	return point, nil
}

// func (dg *DatabaseGame) removePoint() (Point, error) {

// }

func (dg *DatabaseGame) getPointByPlayerId(point PaginationPoint) ([]Point, error) {
	var points []Point
	query := dg.db.Model(&Point{}).Where("player_id = ?", point.PlayerId)
	
	if point.Pagination.SortBy == "" {
		point.Pagination.SortBy = "id"
	}

	if point.Pagination.Sort == "" {
		point.Pagination.Sort = "ASC"
	}

	query = query.Order(
		fmt.Sprintf("%s %s", point.Pagination.SortBy, point.Pagination.Sort)
	).
		Limit(point.Pagination.Size).
		Offset((max(point.Pagination.Page) - 1) * point.Pagination.Size)

	if err := query.Find(&point).Error; err != nil {
		return nil, err
	}

	return points, nil
}
