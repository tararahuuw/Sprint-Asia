import React, {useState, useEffect} from 'react';
import './App.css';
import {Button, Form, ProgressBar, Modal} from 'react-bootstrap';
import SweetAlert from 'react-bootstrap-sweetalert';
import { FaPencilAlt, FaPlus, FaTrash, FaRegCalendarAlt } from "react-icons/fa";
import axios from 'axios';
import { identity } from 'lodash';

function App() {
  // Data State
  const [dataTodo, setDataTodo] = useState([])
  const [dataInput, setDataInput] = useState(true)
  const [dataType, setDataType] = useState()
  const [selectedOption, setSelectedOption] = useState();
  const [dataScout, setDataScout] = useState()
  const [isLoading, setIsLoading] = useState(false);

  // Modals
  const [show, setShow] = useState(false);
  const handleClose = () => setShow(false);
  const handleShow = (id, type) => {
    setDataScout(id)
    setDataType(type)
    setShow(true);
  }

  const handleCheckboxClick = (id, value, type) => {
    modifyData(id, type, value)
  };


  // UseEffect
  useEffect(() => {
    getAllData()
  },[selectedOption])

  useEffect(() => {
  }, [dataTodo])

  // Function
  async function getAllData() {
    setIsLoading(true);
  
    try {
  
      let response;
      if (selectedOption === "Complete") {
        response = await axios.get('http://localhost:8080/task/complete');
        setDataInput(false)
      } else {
        response = await axios.get('http://localhost:8080/task/ongoing');
        setDataInput(true)
      }
  
      const result = response.data;
      console.log("hasil", result);
  
      // Update state with the fetched data
      setDataTodo(result);
      setIsLoading(false);
    } catch (error) {
      console.error('Error fetching data:', error);
      setIsLoading(false);
    }
  }
  

  // Create
  async function modifyData(id, dataType, value) {
    setIsLoading(true);
    try {
      if (dataType == "Create New Task") {
        await axios.post('http://localhost:8080/task/create', {name:dataInput});
      }
      else if (dataType == "Create SubTask") {
        await axios.post('http://localhost:8080/subtask/create', {Id_Task:id, Name:dataInput});
      }
      else if (dataType == "Edit Task") {
        await axios.put(`http://localhost:8080/task/update/${id}`, {Name:dataInput});
        getAllData()
      }
      else if (dataType == "Edit SubTask") {
        await axios.put(`http://localhost:8080/subtask/update/${id}`, {Name:dataInput});
        getAllData()
      }
      else if (dataType == "Edit Deadline") {
        var input = dataInput+":00Z"
        await axios.put(`http://localhost:8080/task/deadline/${id}`, {deadline:input});
      }
      else if (dataType == "Complete Task") {
        console.log("masuk", dataInput)
        var data = await axios.put(`http://localhost:8080/task/complete/${id}`, {complete:value});
        console.log(data)
      }
      else if (dataType == "Complete SubTask") {
        console.log("masuk subtask complete")
        var data = await axios.put(`http://localhost:8080/subtask/complete/${id}`, {complete:value});
        console.log(data)
      }
      else if (dataType == "Delete Task") {
        await axios.delete(`http://localhost:8080/task/delete/${id}`);
      }
      else if (dataType == "Delete SubTask") {
        await axios.delete(`http://localhost:8080/subtask/delete/${id}`);
      }
      console.log("Masuk Modify")
      console.log(id, dataType, dataInput)
    } catch (error) {
      console.log(dataInput)
      console.error('Error fetching data:', error);
    } finally {
      getAllData()
      handleClose()
    }
  }

  const handleDropdownChange = (e) => {
    setSelectedOption(e.target.value);
  };

  function formatDateTime(dateTimeString) {
    const dateTime = new Date(dateTimeString);
  
    // Get date components
    const year = dateTime.getFullYear();
    const month = (dateTime.getMonth() + 1).toString().padStart(2, '0');
    const day = dateTime.getDate().toString().padStart(2, '0');
  
    // Get time components
    const hours = dateTime.getHours().toString().padStart(2, '0');
    const minutes = dateTime.getMinutes().toString().padStart(2, '0');
  
    // Format as "MM-DD-YYYY (HH:mm)"
    const formattedDateTime = `${month}-${day}-${year} (${hours}:${minutes})`;
  
    return formattedDateTime;
  }

  return (
    <div className="App">
      <div className='container-app'>
        <div className='title-app'>
          <p>Todo List Sprint Asia</p>
        </div>
        <div className='header-app'>
          <Button onClick={(e) => handleShow(null, "Create New Task")}>Create New Task</Button>
          <div>
            <Form.Select value={selectedOption} onChange={handleDropdownChange}>
              <option value="On Going">On Going</option>
              <option value="Complete">Complete</option>
            </Form.Select>
          </div>
        </div>
        <div className='content-app'>
          <div className='container-content-app'>
            {
              dataTodo == null ?
                <div className='empty-content-app'>
                  <p>No Todos</p>
                </div>
                :
                <div className="fill-content-app">
                  {
                    dataTodo.map((data) => (
                      <>
                        <div className='container-fill-content-app'>
                          <div className='container-task-content-app'>
                            <div className='container-title-task'>
                              <Form.Check
                                type="checkbox"
                                id={data.Id}
                                checked={data.Complete}
                                onChange={(e) => handleCheckboxClick(data.Id, e.target.checked, "Complete Task")}
                              />
                              <div>
                                {
                                  data.Complete == false ?
                                    <p>{data.Name}</p>
                                    :
                                    <p><strike>{data.Name}</strike></p>
                                }
                                {
                                  data.Deadline != null && data.Deadline != "0001-01-01T00:00:00Z" ?
                                    <p>Deadline : {formatDateTime(data.Deadline)}</p>
                                    :
                                    null
                                }
                              </div>
                            </div>
                            <div className='container-action-task'>
                              {
                                data.Expired == true ?
                                  <div className='container-expired'>
                                    <p>Expired</p>
                                  </div>
                                  :
                                  null
                              }
                              <div className='container-icon-task' onClick={(e) => handleShow(data.Id, "Create SubTask")}>
                                <FaPlus className='icon-task'/>
                              </div>
                              <div className='container-icon-task' onClick={(e) => handleShow(data.Id, "Edit Deadline")}>
                                <FaRegCalendarAlt className='icon-task'/>
                              </div>
                              <div className='container-icon-task' onClick={(e) => handleShow(data.Id, "Edit Task")}>
                                <FaPencilAlt className='icon-task'/>
                              </div>
                              <div className='container-icon-task' onClick={(e) => modifyData(data.Id, "Delete Task")}>
                                <FaTrash className='icon-task'/>
                              </div>
                            </div>
                          </div>
                          {
                            data.SubTasks.length != 0 ?
                              <div>
                                <hr />
                                <div className='header-subtask'>
                                  <p style={{margin:0}}>SubTask :</p>
                                  <ProgressBar animated now={data.Progress} label={`${data.Progress}%`} className="progressbar-task" />
                                </div>
                                {
                                  data.SubTasks.map((child) => (
                                    <>
                                      <div className='container-task-content-app'>
                                        <div className='container-title-task'>
                                          <Form.Check
                                            type="checkbox"
                                            id={child.Id}
                                            checked={child.Complete}
                                            onChange={(e) => handleCheckboxClick(child.Id, e.target.checked, "Complete SubTask")}
                                          />
                                          {
                                            child.Complete == false ?
                                              <p>{child.Name}</p>
                                              :
                                              <p><strike>{child.Name}</strike></p>
                                          }
                                        </div>
                                        <div className='container-action-task'>
                                          <div className='container-icon-task' onClick={(e) => handleShow(child.Id, "Edit SubTask")}>
                                            <FaPencilAlt className='icon-task'/>
                                          </div>
                                          <div className='container-icon-task' onClick={(e) => modifyData(child.Id, "Delete SubTask")}>
                                            <FaTrash className='icon-task'/>
                                          </div>
                                        </div>
                                      </div>
                                    </>
                                  ))
                                }
                              </div>
                              :
                              null
                          }
                        </div> 
                      </>
                    ))
                  }
                </div>
            }
          </div>
        </div>
      </div>
      <Modal show={show} centered>
        <Modal.Header >
          <Modal.Title>{dataType}</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          {
            dataType == "Edit Deadline" ?
              <Form>
                <Form.Group className="mb-3" controlId="exampleForm.ControlInput1">
                  <Form.Label>Deadline</Form.Label>
                  <Form.Control
                    type="datetime-local"
                    autoFocus
                    onChange={(e) => setDataInput(e.target.value)}
                  />
                </Form.Group>
              </Form>
              :
              <Form>
                <Form.Group className="mb-3" controlId="exampleForm.ControlInput2">
                  <Form.Label>Name</Form.Label>
                  <Form.Control
                    type="text"
                    autoFocus
                    onChange={(e) => setDataInput(e.target.value)}
                  />
                </Form.Group>
              </Form>
          }
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={handleClose}>
            Close
          </Button>
          <Button variant="primary" onClick={e => modifyData(dataScout, dataType)}>
            Save Changes
          </Button>
        </Modal.Footer>
      </Modal>
      {isLoading ? (
        <SweetAlert title="Loading..." showConfirm={false} />
      ) : null}
    </div>
  );
}

export default App;