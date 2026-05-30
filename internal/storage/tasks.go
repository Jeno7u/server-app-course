package storage
package storage

import (
	"sync"
	"github.com/Jeno7u/server-app-course/internal/models"
)

type TaskStorage struct {
	sync.Mutex
	tasks  map[int]models.Task
	nextID int
}

func NewTaskStorage() *TaskStorage {
	return &TaskStorage{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}

func (s *TaskStorage) Create(task models.Task) models.Task {
	s.Lock()
	defer s.Unlock()
	task.ID = s.nextID
	s.tasks[s.nextID] = task
	s.nextID++
	return task
}

func (s *TaskStorage) GetByID(id int) (models.Task, bool) {
	s.Lock()
	defer s.Unlock()
	task, ok := s.tasks[id]
	return task, ok
}

func (s *TaskStorage) GetAll() []models.Task {
	s.Lock()
	defer s.Unlock()
	res := make([]models.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		res = append(res, t)
	}
	return res
}

func (s *TaskStorage) Update(task models.Task) {
	s.Lock()
	defer s.Unlock()
	s.tasks[task.ID] = task
}

func (s *TaskStorage) Delete(id int) bool {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.tasks[id]; ok {
		delete(s.tasks, id)
		return true
	}
	return false
}

func (s *TaskStorage) Reset() {
	s.Lock()
	defer s.Unlock()
	s.tasks = make(map[int]models.Task)
	s.nextID = 1
}

var GlobalTaskStorage = NewTaskStorage()
