import { CallOptions } from 'grpc-web'
import Navigo from 'navigo'
import { AuthServiceClient, LoginRequest, AuthUserRequest } from './proto/services_grpc_web_pb'

const router = new Navigo()
const authClient = new AuthServiceClient('http://localhost:9001')

router
    .on("/", function() {
        document.body.innerHTML = "Home"
    })
    .on("/login", function() {
        document.body.innerHTML = ""
        const loginDiv = document.createElement('div')
        loginDiv.classList.add('login-div')

        const loginLabel = document.createElement('h1')
        loginLabel.innerText = "Login"
        loginDiv.appendChild(loginLabel)

        const loginForm = document.createElement('form')

        const loginInputLabel = document.createElement('label')
        loginInputLabel.classList.add('input-label')
        loginInputLabel.innerText = "Username/Email"
        loginInputLabel.setAttribute('for', 'login-input')
        loginForm.appendChild(loginInputLabel)

        const loginInput = document.createElement('input')
        loginInput.id = "login-input"
        loginInput.setAttribute('type', 'text')
        loginInput.setAttribute('placeholder', 'hitesh5678')
        loginForm.appendChild(loginInput)

        const passwordInputLabel = document.createElement('label')
        passwordInputLabel.classList.add('input-label')
        passwordInputLabel.innerText = "Password"
        passwordInputLabel.setAttribute('for', 'password-input')
        loginForm.appendChild(passwordInputLabel)

        const passwordInput = document.createElement('input')
        passwordInput.id = "password-input"
        passwordInput.setAttribute('type', 'password')
        passwordInput.setAttribute('placeholder', 'set-complex-password')
        loginForm.appendChild(passwordInput)

        const submitBtn = document.createElement('button')
        submitBtn.innerText = "Login"
        loginForm.appendChild(submitBtn)

        loginForm.addEventListener('submit', event => {
            let i = 0
            event.preventDefault()
            if (i != 0) return
            ++i
            let req = new LoginRequest()
            req.setLogin(loginInput.value)
            req.setPassword(passwordInput.value)
            authClient.login(req, {}, (err, res) => {
                if (err) return alert(err.message)
                localStorage.setItem('token', res.getToken())
                req = new AuthUserRequest()
                req.setToken(res.getToken())
                let j = 0
                authClient.authUser(req, {}, (err, res) => {
                    if (j != 0) return
                    ++j
                    if (err) return alert(err.message)
                    const user = { id: res.getId(), email: res.getEmail(), username: res.getUsername() }
                    localStorage.setItem('user', JSON.stringify(user))
                })
            })
        })

        loginDiv.appendChild(loginForm)

        document.body.appendChild(loginDiv)
    })
    .resolve()