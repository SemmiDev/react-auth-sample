import { useState } from 'react';

function useAuthDetails() {
  const getAuthDetail = () => {
    return JSON.parse(localStorage.getItem('authDetails'));
  };

  const [authDetails, setAuthDetails] = useState(getAuthDetail());
  console.log(authDetails);

  const saveAuthDetails = userDetail => {
    localStorage.setItem('authDetails', JSON.stringify(userDetail));
    setAuthDetails(userDetail);
  };

  return {
    setAuthDetail: saveAuthDetails,
    authDetails
  }
}

export default useAuthDetails;