import {router,authClient} from '../app'
import { AuthUserRequest, SignupRequest, UsernameUsedRequest, EmailUsedRequest } from '../proto/services_grpc_web_pb'

function validateEmail(email) {
    const re = /^(([^<>()[\]\.,;:\s@\"]+(\.[^<>()[\]\.,;:\s@\"]+)*)|(\".+\"))@(([^<>()[\]\.,;:\s@\"]+\.)+[^<>()[\]\.,;:\s@\"]{2,})$/i;
    return re.test(String(email).toLowerCase());
}

function Signup() {
    document.body.innerHTML = ""

    const signupDiv = document.createElement('div')
    signupDiv.classList.add('auth-div')

    const signupLabel = document.createElement('h1')
    signupLabel.innerText = "Signup"
    signupDiv.appendChild(signupLabel)

    const signupForm = document.createElement('form')

    const usernameInputLabel = document.createElement('label')
    usernameInputLabel.classList.add('input-label')
    usernameInputLabel.innerText = "Username"
    usernameInputLabel.setAttribute('for', 'username-input')
    signupForm.appendChild(usernameInputLabel)

    const usernameInput = document.createElement('input')
    usernameInput.id = "username-input"
    usernameInput.setAttribute('type', 'text')
    usernameInput.setAttribute('placeholder', 'Hitesh5678')
    signupForm.appendChild(usernameInput)

    usernameInput.addEventListener("input", () => {
        usernameError.innerText = ""
        const username = usernameInput.value
        if (username.length < 4 || username.length > 20) {
            usernameError.innerText = "User name must be atleast 4 characters long and no more than 20 characters."
            return
        }
        let req = new UsernameUsedRequest()
        req.setUsername(usernameInput.value)
        authClient.usernameUsed(req, {}, (err, res) => {
            if (err) usernameError.innerText = err.message
            if (res.getUsed()) usernameError.innerText = "Username already in use."
        })
    })

    const usernameError = document.createElement('div')
    usernameError.id = "username-error"
    usernameError.classList.add('error')
    signupForm.appendChild(usernameError)

    const emailInputLabel = document.createElement('label')
    emailInputLabel.classList.add('input-label')
    emailInputLabel.innerText = "Email"
    emailInputLabel.setAttribute('for', 'email-input')
    signupForm.appendChild(emailInputLabel)

    const emailInput = document.createElement('input')
    emailInput.id = "email-input"
    emailInput.setAttribute('type', 'email')
    emailInput.setAttribute('placeholder', 'hitesh@gmail.com')
    signupForm.appendChild(emailInput)

    emailInput.addEventListener("input", () => {
        emailError.innerText = ""
        const email = emailInput.value
        if (email.length < 7 || email.length > 35 || !validateEmail(email)) {
            emailError.innerText = "Email ID not correct."
            return
        }
        let req = new EmailUsedRequest()
        req.setEmail(emailInput.value)
        authClient.emailUsed(req, {}, (err, res) => {
            if (err) emailError.innerText = err.message
            if (res.getUsed()) emailError.innerText = "Email ID already in use."
        })
    })

    const emailError = document.createElement('div')
    emailError.id = "email-error"
    emailError.classList.add('error')
    signupForm.appendChild(emailError)

    const passwordLabel = document.createElement('label')
    passwordLabel.classList.add('input-label')
    passwordLabel.innerText = "Password"
    passwordLabel.setAttribute('for', 'password-input')
    signupForm.appendChild(passwordLabel)

    const passwordInput = document.createElement('input')
    passwordInput.id = "password-input"
    passwordInput.setAttribute('type', 'password')
    passwordInput.setAttribute('placeholder', 'password')
    signupForm.appendChild(passwordInput)

    passwordInput.addEventListener("input", () => {
        passwordError.innerText = ""
        const password = passwordInput.value
        if (password.length < 8 || password.length > 120) {
            passwordError.innerText = "User name must be atleast 8 characters long and no more than 120 characters."
        }
    })

    const passwordError = document.createElement('div')
    passwordError.id = "password-error"
    passwordError.classList.add('error')
    signupForm.appendChild(passwordError)

    const signupBtn = document.createElement('button')
    signupBtn.innerText = "Signup"
    signupForm.appendChild(signupBtn)

    signupForm.addEventListener('submit', event => {
        let i = 0
        event.preventDefault()
        if (i != 0) return
        ++i
        
        if (usernameInput.value == "" || usernameError.innerText != "" || emailInput.value == "" || emailError.innerText != "" || passwordInput.value == "" || passwordError.innerText != "") return
        
        let req = new SignupRequest()
        req.setUsername(usernameInput.value)
        req.setEmail(emailInput.value)
        req.setPassword(passwordInput.value)

        authClient.signup(req, {}, (err, res) => {
            if (err) return alert(err.message)
            localStorage.setItem("token", res.getToken())

            req = new AuthUserRequest()
            req.setToken(res.getToken())

            authClient.authUser(req, {}, (err, res) => {
                if (err) return alert(err.message)
                const user = { id: res.getId(), email: res.getEmail(), username: res.getUsername() }
                localStorage.setItem('user', JSON.stringify(user))
                router.navigate('/')  
            })
        })
    })

    const loginText = document.createElement("label")
    loginText.classList.add('input-label')
    loginText.innerText = "Already registered? Please login."
    loginText.setAttribute('for', 'login-btn')
    signupForm.appendChild(loginText)

    const loginBtn = document.createElement('button')
    loginBtn.id = "login-btn"
    loginBtn.innerText = "Login"
    loginBtn.addEventListener("click", () => {
        router.navigate('/login')
    })
    signupForm.appendChild(loginBtn)

    signupDiv.appendChild(signupForm)

    document.body.appendChild(signupDiv)

}

export {Signup}