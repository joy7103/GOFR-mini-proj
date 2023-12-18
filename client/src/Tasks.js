import { Component } from "react";
import {
    addTask,
    getTasks,
    updateTask,
    deleteTask,
} from "./services/taskServices";

class Tasks extends Component {
    state = { tasks: [], currentTask: "" };

    async componentDidMount() {
        try {
            const { data } = await getTasks();

            console.log("Received tasks:", data);

            this.setState({ tasks: data });
        } catch (error) {
            console.log(error);
        }
    }

    handleChange = ({ currentTarget: input }) => {
        this.setState({ currentTask: input.value });
    };

    handleSubmit = async (e) => {
        e.preventDefault();
        const originalTasks = this.state.tasks || [];
        try {
            const { data } = await addTask({ task: this.state.currentTask });
            console.log("Received data from addTask:", data);
            const tasks = [...originalTasks, data]; // Create a new array with the updated data
            this.setState({ tasks, currentTask: "" });
        } catch (error) {
            console.log(error);
        }
    };

    handleUpdate = async (currentTask) => {
        const originalTasks = [...this.state.tasks];
        try {
            const tasks = [...originalTasks];
            const index = tasks.findIndex((task) => task._id === currentTask);
            tasks[index] = { ...tasks[index] };
            tasks[index].completed = !tasks[index].completed;
            
            await updateTask(currentTask, {
                completed: tasks[index].completed,
            });
            setTimeout(() => {
                window.location.reload();
              }, 5);

            this.setState({ tasks }, () => {
                console.log("State updated after handleUpdate:", this.state.tasks);
            });

            
            
        } catch (error) {
            this.setState({ tasks: originalTasks });
            console.log(error);
        }
    };

    handleDelete = async (currentTask) => {
        // Log the received currentTask to understand its structure
        console.log("Received currentTask for deletion:", currentTask);
    
        const originalTasks = this.state.tasks;
    
        try {
            // Ensure that the task ID is not undefined before proceeding
            if (!currentTask) {
                console.error("Task ID is undefined. Aborting deletion.");
                return;
            }
    
            const tasks = originalTasks.filter((task) => task._id !== currentTask);
            this.setState({ tasks });


            await deleteTask(currentTask); // Pass the task ID directly as a string

            setTimeout(() => {
                window.location.reload();
              }, 5);
            

        } catch (error) {
            this.setState({ tasks: originalTasks });
            console.log("Error deleting task:", error);
        }
    };
}

export default Tasks;
