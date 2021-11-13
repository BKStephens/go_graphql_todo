import axios from 'axios';

const register = (username, email, password) => {
  return axios.post('/signup', {
    username,
    email,
    password,
  });
};

const login = (username, password) => {
  return axios
    .post('/login', {
      username,
      password,
    })
    .then((response) => {
      if (response.data.token) {
        localStorage.setItem('user', JSON.stringify(response.data));
      }

      return response.data;
    });
};

const logout = () => {
  localStorage.removeItem('user');
};

const getCurrentUser = () => {
  return JSON.parse(localStorage.getItem('user'));
};

const exports = {
  register,
  login,
  logout,
  getCurrentUser,
};
export default exports;
