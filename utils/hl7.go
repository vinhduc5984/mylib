package utils

import (
	"log"
	"os"
	"strings"
)

func WriteHl7OrderRis(path string, indicationItemId, patientCode, patientName string, patientBirthDate int64, patientGender int32, patientAddress, risCode, risName, departmentName, indicatorName, orderControl string, orderDateTime int64, relevantClinicalInformation, accessionNo, studyUID, modality, room, aet, sending, receiving string) (string, error) {
	hl7Message, err := getHl7OrderRisMessage(indicationItemId, patientCode, patientName, patientBirthDate, patientGender, patientAddress, risCode, risName, departmentName, indicatorName, orderControl, orderDateTime, relevantClinicalInformation, accessionNo, studyUID, modality, room, aet, sending, receiving)
	if err != nil {
		return "", err
	}
	existed, err := IsFileExisted(path)
	if err != nil {
		return "", err
	}
	if !existed {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Println(err)
			return "", err
		}
	}
	hl7FilePath := path + "/" + accessionNo + ".hl7"
	fileExisted, err := IsFileExisted(hl7FilePath)
	if err != nil {
		return "", err
	}
	if fileExisted {
		if err = DeleteFile(hl7FilePath); err != nil {
			return "", err
		}
	}

	if err = CreateFile(hl7FilePath); err != nil {
		return "", err
	}
	if err = WriteFile(hl7FilePath, hl7Message); err != nil {
		return "", err
	}

	return hl7Message, nil
}

func getHl7OrderRisMessage(indicationItemId string, patientCode string, patientName string, patientBirthDate int64,
	patientGender int32, patientAddress string, risCode string, risName, departmentName,
	indicatorName, orderControl string, orderDateTime int64, relevantClinicalInformation, accessionNo, studyUID,
	modality, room, aet, sending, receiving string) (string, error) {
	gender := "F"
	if patientGender == 1 {
		gender = "M"
	}
	orderDate := FormatDateWithLayout(orderDateTime, "20060102150405") //yyyyMMddHHmmss
	birthDate := FormatDateWithLayout(patientBirthDate, "20060102")    //yyyyMMdd

	risCode = strings.TrimSpace(risCode)
	risName = strings.TrimSpace(risName)
	itemCodeName := risCode + "^" + risName

	hl7 := segMSH(orderDate, sending, receiving)
	hl7 += "\r\n" + segEVN(indicationItemId)
	hl7 += "\r\n" + segPID(patientCode, patientName, birthDate, gender, patientAddress)
	hl7 += "\r\n" + segPV1(departmentName, indicatorName)
	hl7 += "\r\n" + segORC(orderControl, orderDate)
	hl7 += "\r\n" + segOBR(indicationItemId, itemCodeName, orderDate, relevantClinicalInformation, accessionNo, room, aet, modality)
	//obx
	hl7 += "\r\n" + segZDS(studyUID)

	return hl7, nil
}

func segMSH(orderDate, sending, receiving string) string {
	return "MSH|^~\\&|" + sending + "|" + sending + "|" + receiving + "|" + receiving + "|" + orderDate + "||ORM^O01||P||2.3.1|||||||"
}

func segEVN(orderDate string) string {
	return "EVN|O01|" + orderDate
}

func segPID(patientCode, patientName, patientBirthDate, patientGender, patientAddress string) string {
	return "PID|1||" + patientCode + "||" + patientName + "||" + patientBirthDate + "|" + patientGender + "|||" + patientAddress + "|||||||||||||||||||"
}

func segPV1(departmentName, indicatorName string) string {
	return "PV1|||" + departmentName + "|||||" + indicatorName + "||||||||||||||||||||||||||||||||||||||||||||"
}

func segORC(orderControl, orderDate string) string {
	return "ORC|" + orderControl + "||||IP||" + orderDate + "||||||||||||"
}

// relevantClinicalInformation: chẩn đoán
func segOBR(indicationItemId, itemCodeName, orderDate, relevantClinicalInformation, accessionNo, room, eat, modality string) string {
	return "OBR||" + indicationItemId + "|" + indicationItemId + "|" + itemCodeName + "||" + orderDate + "|||||||" + relevantClinicalInformation + "|||||" + accessionNo + "|" + room + "|" + eat + "|" + modality + "|SCHEDULED||" + modality + "|||||||||||||||||||"
}

func segZDS(studyUID string) string {
	return "ZDS|" + studyUID
}
