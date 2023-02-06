import { useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '../../contexts/Auth';

export function LoginPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const auth = useAuth();

  return (
    <div className='dark2'>
      <div className="center">
        <div className="login-form-wrapper">
          <h1>Login</h1>

          <form
            method="post"
            onSubmit={async (event) => {
              event.preventDefault();

              const formData = new FormData(event.currentTarget);
              const email = formData.get('email');
              const password = formData.get('password');

              await auth.login({ email, password });
              const from = location.state?.from?.pathname ?? '/';
              navigate(from, { replace: true });
            }}
          >
            <div className="form-element">
              <label>Email</label>
              <input
                className="input_filed"
                name="email"
                type="email"
                autoComplete="off"
                placeholder=" "
                required
              />
              <span></span>
            </div>
            <div className="form-element">
              <label>Password</label>
              <input
                className="input_filed"
                name="password"
                type="password"
                autoComplete="off"
                placeholder=" "
                required
              />
              <span></span>
            </div>
            <input type="submit" value='Login' />
          </form>
        </div>
      </div>
    </div>
  );
}
