import {router} from '../app'

function Home() {
    document.body.innerHTML = ""

    const homeDiv = document.createElement('div')
    homeDiv.classList.add('home-div')

    const user = JSON.parse(localStorage.getItem('user'))
    console.log(user)
    if (user == null) {

        const buttonContainer = document.createElement('div')
        buttonContainer.classList.add('button-container')

        const loginBtn = document.createElement('button')
        loginBtn.innerText = "Login"
        loginBtn.addEventListener("click", () => {
            router.navigate('/login')
        })

        buttonContainer.appendChild(loginBtn)
        
        const text = document.createElement('h2')
        text.innerText = "You are not authenticated!!"

        buttonContainer.appendChild(text)

        const signupBtn = document.createElement('button')
        signupBtn.innerText = "Signup"
        signupBtn.addEventListener("click", () => {
            router.navigate('/signup')
        })

        buttonContainer.appendChild(signupBtn)
        homeDiv.appendChild(buttonContainer)

    } else {

        const authText = document.createElement('div')
        authText.classList.add('auth-text')
        authText.innerText = `You are logged in as ${user.username}.\nYour email is ${user.email}.`
        homeDiv.appendChild(authText)

        const logoutBtn = document.createElement('button')
        logoutBtn.innerText = "Logout"
        logoutBtn.addEventListener("click", () => {
            localStorage.setItem('user', null)
            localStorage.setItem('token', null)
            router.navigate('/login')
        })

        homeDiv.appendChild(logoutBtn)
    }

    document.body.appendChild(homeDiv)

}

export {Home}