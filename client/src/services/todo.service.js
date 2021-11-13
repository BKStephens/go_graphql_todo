import axios from 'axios';
import authHeader from './auth-header';

const getTodos = () => {
  return axios.get('/graphql?query={todoList{id,text,done}}', { headers: authHeader() });
};

const addTodo = (text) => {
  const encodedText = encodeURIComponent(text)
  return axios.get(`/graphql?query=mutation+_{createTodo(text:"${encodedText}"){id,text,done}}`, { headers: authHeader() });
};

const updateTodo = (id, done) => {
  return axios.get(`/graphql?query=mutation+_{updateTodo(id:${id},done:${done}){id,text,done}}`, { headers: authHeader() });
}

const exports = {
  getTodos,
  addTodo,
  updateTodo,
};
export default exports;
