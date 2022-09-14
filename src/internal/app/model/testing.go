package model

import "testing"

func TestPatient(t *testing.T) *Patient {
	t.Helper()

	return &Patient{
		Name: "Полина",
		Surname: "Сироткина",
		Gender: "Женский",
		Birth_year: 2001,
		Phone: "89998887766",
		Email: "polina@example.ru",
		Password: "123456",
	}
}

func TestDoctor(t *testing.T) *Doctor {
	t.Helper()

	return &Doctor{
		Name: "Полина",
		Surname: "Сироткина",
		Work_since: 2000,
		Spec_id: 1,
		Email: "polina@example.ru",
		Password: "123456",
	}
}

func TestAdmin(t *testing.T) *Admin {
	t.Helper()

	return &Admin{
		Name: "Полина",
		Surname: "Сироткина",
		Email: "polina@example.ru",
		Password: "123456",
	}
}

func TestSpecialization(t *testing.T) *Specialization {
	t.Helper()

	return &Specialization{
		Name: "Хирург",
		Salary: 20000,
	}
}

func TestMedicine(t *testing.T) *Medicine {
	t.Helper()

	return &Medicine{
		Name: "Ношпа",
		Cost: 100,
	}
}

func TestRecord(t *testing.T) *Record {
	t.Helper()

	return &Record{
		Visit_id: 1,
		Disease_id: 1,
		Medicine_id: 1,
	}
}

func TestDisease(t *testing.T) *Disease {
	t.Helper()

	return &Disease{
		Name: "Насморк",
		Spec_id: 1,
	}
}

func TestVisit(t *testing.T) *Visit {
	t.Helper()

	return &Visit{
		Status: "Active",
		Doctor_id: 1,
		Patient_id: 1,
	}
}