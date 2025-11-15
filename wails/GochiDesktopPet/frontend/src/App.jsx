import React, { useState } from 'react';
//import { GochiDesktopPet } from '../../wailsjs/go/main/App.js';
import './style.css';

export default function App() {
  const [tasks, setTasks] = useState([]);
  const [mood, setMood] = useState('Neutral.');

  const addTask = async () => {
    const task = prompt('New task:');
    if (task) {
      await GochiDesktopPet.AddTask(task);
      setTasks([...tasks, task]);
      const newMood = await GochiDesktopPet.Mood();
      setMood(newMood);
    }
  };

  return (
    <div id="pet">
      <h1 style={{fontFamily: "Nougat-ExtraBlack"}}>Gochi!</h1>
      <p style={{fontFamily: "AlteHaasGroteskRegular"}}>Current mood: {mood}</p>
      <button onClick={addTask}>Add Task</button>
      <ul>
        {tasks.map((t, i) => (
          <li key={i}>{t}</li>
        ))}
      </ul>
    </div>
  );
}
