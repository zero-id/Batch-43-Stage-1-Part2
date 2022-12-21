function getData() {

    let name = document.getElementById(`name`).value
    let email = document.getElementById(`email`).value
    let phoneNumber = document.getElementById(`telp`).value
    let subject = document.getElementById(`subject`).value
    let message = document.getElementById(`message`).value

    let emailReciver =`latiff@gmail.com`

    let mailTo = document.createElement(`a`)
    mailTo.href = `mailto:${emailReciver}?subject=${subject}&body=Hallo nama saya ${name}, ${message}, nomor telpon saya ${phoneNumber}, email saya ${email}`

    let users = {
        myName: name,
        myEmail: email,
        myPhone: phoneNumber,
        mySubject: subject,
        myMessage: message
    }

    if (name === ``) {
        return alert(`Nama harus di isi!`)
    } else if (email === ``) {
        return alert(`Email harus di isi!`)
    } else if (phoneNumber ===  ``) {
        return alert(`Phone Number harus di isi!`)
    } else if (subject === ``) {
        return alert(`Subject harus di pilih!`)
    } else if (message === ``) {
        return alert(`message harus di isi!`)
    } else {
        console.log(users), mailTo.click()
    }

}







