function checkInput() {
    const cBox = document.getElementById('more-Boxes');
    const addForm = document.getElementById('open-form-1');
    const clickOnSecondBtn = document.getElementById('clientTwo');
    const addThirdForm = document.getElementById('open-form-3');
    const clickOnThirdBtn = document.getElementById('clientThree');
    const addFourthForm = document.getElementById('open-form-4');

    if (cBox.checked == true) {
        addForm.style.display = 'block'
    } else {
        addForm.style.display = 'none'
    }

    if (clickOnSecondBtn.checked == true) {
        addThirdForm.style.display = 'block'
    } else {
        addThirdForm.style.display = 'none'
    }

    if (clickOnThirdBtn.checked == true) {
        addFourthForm.style.display = 'block'
    } else {
        addFourthForm.style.display = 'none'
    }
}
