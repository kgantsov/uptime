import React from 'react';

export interface CurrentUser {
  email: string;
  first_name: string;
  last_name: string;
}

export const CurrentUserContext = React.createContext<
  [CurrentUser | null, (currentUser: CurrentUser) => void]
>([null, (...any: any): void => void 0]);

export const useCurrentUser = () => React.useContext(CurrentUserContext);

export const CurrentUserProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const value = React.useState<CurrentUser | null>(null);

  return (
    <CurrentUserContext.Provider value={value}>
      {children}
    </CurrentUserContext.Provider>
  );
};
