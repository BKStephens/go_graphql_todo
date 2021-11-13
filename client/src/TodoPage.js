import React, {useEffect, useState} from 'react';

import TodoList from './TodoList';
import TodoForm from './TodoForm';
import TodoService from './services/todo.service';

const TodoPage = () => {
  const [todos, setTodos] = useState([]);
  useEffect(() => {
    TodoService.getTodos().then(
      (response) => setTodos(response.data.data.todoList),
      (error) => {
        console.log(error);
      }
    );
  }, []);

  const addTodo = (text) => {
    TodoService.addTodo(text).then(
      (response) => {
        setTodos(todos.concat(response.data.data.createTodo));
      },
      (error) => {
        console.log(error);
      }
    )
  };

  const toggleTodo = (id) => {
    const todo = todos.find(x => x.id === id);
    TodoService.updateTodo(id, !todo.done).then(
      (response) => {
        setTodos(todos.map(x => x.id === id ? response.data.data.updateTodo : x))
      },
      (error) => {
        console.log(error);
      }
    )
  };

  return (
    <div>
      <h1>Todos</h1>
      <TodoList todos={todos} handleToggle={toggleTodo} />
      <TodoForm addTodo={addTodo} />
    </div>
  );
};

export default TodoPage;
