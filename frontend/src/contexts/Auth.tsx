import React, { useCallback } from 'react';
import { Navigate, useLocation } from 'react-router-dom';

export const authProvider = {
  async login(userCredentials: UserCredentials) {
    try {
      const resp = await fetch('/API/v1/tokens', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(userCredentials),
      });

      if (!resp.ok) {
        return;
      }

      return await resp.json();
    } catch {}
  },

  async logout() {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };

    const token = localStorage.getItem('token');

    if (token !== null) {
      headers.Authorization = `Bearer ${token}`;
    }

    try {
      const resp = await fetch('/API/v1/tokens', {
        method: 'DELETE',
        headers: headers,
      });

      if (!resp.ok) {
        console.log('Failed to remove token. Got status: ', resp.status);
      }

      return true;
    } catch (error) {
      return false;
    }
  },
};

interface UserCredentials {
  email: string;
  password: string;
}

export interface AuthContextProps {
  userCredentials: UserCredentials | null;
  token: string | null;
  login: (userCredentials: UserCredentials) => void;
  logout: () => void;
}

export const AuthContext = React.createContext<AuthContextProps>({
  userCredentials: null,
  token: null,
  login: (...any: any): void => void 0,
  logout: (...any: any): void => void 0,
});

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [userCredentials, setUserCredentials] =
    React.useState<UserCredentials | null>(null);
  const [token, setToken] = React.useState<string | null>(
    localStorage.getItem('token'),
  );

  const login = useCallback(async (userCredentials: UserCredentials) => {
    const data = await authProvider.login(userCredentials);

    if (data) {
      const token = data.token;
      localStorage.setItem('token', token);
      setToken(token);
      setUserCredentials(userCredentials);
    }
  }, []);

  const logout = useCallback(async () => {
    const loggedOut = await authProvider.logout();

    if (loggedOut) {
      localStorage.removeItem('token');
      setUserCredentials(null);
      setToken(null);
    }
  }, []);

  return (
    <AuthContext.Provider value={{ userCredentials, token, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  return React.useContext(AuthContext);
};

export const RequireAuth: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const auth = useAuth();
  const location = useLocation();

  if (!auth?.token) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  return <>{children}</>;
};
