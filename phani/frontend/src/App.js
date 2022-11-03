import "./App.css";
// import axios from "axios";
import React, { useState, useEffect } from 'react';
import axios from 'axios';

function App() {
  // var owner_files, shared_files;
  const [file, setFile] = useState()
  const [username, setUsername] = useState("Dropdown")
  const [openDropdown, setOpenDropdown] = useState(false)
  const [selected, setselected] = useState({
    'ownerName': '',
    'fileName': ''
  });
  function handleChange(event) {
    setFile(event.target.files[0])
  }
  // function handleResponse(response) {
  //   response.blob().then(blob => {
  //     const link = document.createElement('a');
  //     const url = URL.createObjectURL(blob);
  //     link.href = url;
  //     link.download = "1.pdf";
  //     link.closest();
  //   });
  // }
  function handleUsernameChange(e) {
    setUsername(e.target.innerText);
    getdata();
    getdata();
    // console.log(e.target.innerText)
  }

  function handleDropdownOpen(e) {
    setOpenDropdown(!openDropdown);
  }

  function handleOwnerSelect(e) {
    console.log(e.target.innerText);
    // setselected()
  }
  function handleSharedSelect(e) {
    console.log(e);
  }
  async function deleteFile(event) {
    event.preventDefault();

    const url = 'http://localhost:8080/delete';
    const formData = new FormData();
    formData.append('ownerName', selected.ownerName);
    formData.append('fileName', selected.fileName);
    await axios.post(url, formData,
      {
        'content-type': 'multipart/form-data',
        'Access-Control-Allow-Origin': '*',
        'Access-Control-Allow-Methods': 'GET, POST, PATCH, PUT, DELETE, OPTIONS',
        'Access-Control-Allow-Headers': 'Origin, Content-Type, X-Auth-Token',
        'crossorigin': true,
        'Access-Control-Allow-Credentials': true,
      }
    ).then((response) => {
      console.log("indelete");
      console.log(response.data);
    });
  }
  async function downloadFile(event) {
    event.preventDefault();

    const url = 'http://localhost:8080/download';
    const formData = new FormData();
    formData.append('ownerName', selected.ownerName);
    formData.append('fileName', selected.fileName);
    await axios.post(url, formData,
      {
        'content-type': 'multipart/form-data',
        'Access-Control-Allow-Origin': '*',
        'Access-Control-Allow-Methods': 'GET, POST, PATCH, PUT, DELETE, OPTIONS',
        'Access-Control-Allow-Headers': 'Origin, Content-Type, X-Auth-Token',
        'crossorigin': true,
        'Access-Control-Allow-Credentials': true,
        'responseType': 'blob',
        'time': 'sting',
        'time_units': 'stin'
      }
    ).then((response) => {
      console.log("indownload");
      console.log(response.data);
      console.log(response);

      // var binaryString = response.data;
      // var binaryLen = binaryString.length;
      // var bytes = new Uint8Array(binaryLen);
      // for (let i = 0; i < binaryLen; i++) {
      //   var ascii = binaryString.charCodeAt(i);
      //   bytes[i] = ascii;
      // }
      // var bytes = new Uint8Array(response.data);
      var data = new Blob([response.data]);

      var link = document.createElement('a');
      link.href = window.URL.createObjectURL(data);
      var fileName = selected.fileName;
      link.download = fileName;
      link.click();
    });
  }
  async function ShareFile(event) {
    event.preventDefault();

    const url = 'http://localhost:8080/share';
    const formData = new FormData();
    var shareName = document.getElementById("shareNameId").value;
    formData.append('ownerName', selected.ownerName);
    formData.append('shareName', shareName);
    formData.append('fileName', selected.fileName);
    await axios.post(url, formData,
      {
        'content-type': 'multipart/form-data',
        'Access-Control-Allow-Origin': '*',
        'Access-Control-Allow-Methods': 'GET, POST, PATCH, PUT, DELETE, OPTIONS',
        'Access-Control-Allow-Headers': 'Origin, Content-Type, X-Auth-Token',
        'crossorigin': true,
        'Access-Control-Allow-Credentials': true,
      }
    ).then((response) => {
      console.log("inshare");
      console.log(response.data);
    });
  }
  const [owner_files, setowner_files] = useState([]);
  const [shared_files, setshared_files] = useState([]);
  useEffect(() => {
    async function getdata(event) {
      const url = 'http://localhost:8080/getData';
      const formData = new FormData();
      formData.append('userName', username);
      await axios.post(url, formData,
        {
          'content-type': 'multipart/form-data',
          'Access-Control-Allow-Origin': '*',
          'Access-Control-Allow-Methods': 'GET, POST, PATCH, PUT, DELETE, OPTIONS',
          'Access-Control-Allow-Headers': 'Origin, Content-Type, X-Auth-Token',
          'crossorigin': true,
          'Access-Control-Allow-Credentials': true,
        }
      ).then((response) => {

        console.log(response.data["owner_files"]);
        setowner_files(response.data["owner_files"]);
        console.log(response.data["shared_files"]);
        setshared_files(response.data["shared_files"]);
      });
    }
    getdata();
  }, [])

  async function getdata(event) {
    const url = 'http://localhost:8080/getData';
    const formData = new FormData();
    formData.append('userName', username);
    await axios.post(url, formData,
      {
        'content-type': 'multipart/form-data',
        'Access-Control-Allow-Origin': '*',
        'Access-Control-Allow-Methods': 'GET, POST, PATCH, PUT, DELETE, OPTIONS',
        'Access-Control-Allow-Headers': 'Origin, Content-Type, X-Auth-Token',
        'crossorigin': true,
        'Access-Control-Allow-Credentials': true,
      }
    ).then((response) => {

      console.log("username", username)
      console.log(response.data["owner_files"]);
      setowner_files(response.data["owner_files"]);
      console.log(response.data["shared_files"]);
      setshared_files(response.data["shared_files"]);
    });
  }

  async function handleSubmit(event) {
    event.preventDefault();
    // const file = event.target.files[0];

    const url = 'http://localhost:8080/upload';
    const formData = new FormData();

    formData.append('file', file);
    formData.append('fileName', file.name);
    formData.append('userName', username);
    const replication = document.getElementById("replication").value;
    formData.append('replication', replication);
    // const config = {
    //   headers: {
    //     'content-type': 'multipart/form-data',
    //     'Access-Control-Allow-Origin': '*',
    //     'Access-Control-Allow-Methods': 'GET, POST, PATCH, PUT, DELETE, OPTIONS',
    //     'Access-Control-Allow-Headers': 'Origin, Content-Type, X-Auth-Token',
    //     'crossorigin': true,
    //     'Access-Control-Allow-Credentials': true,
    //   },
    // };
    console.log(url);
    console.log(formData.data);

    await axios.post(url, formData,
      {
        'content-type': 'multipart/form-data',
        'Access-Control-Allow-Origin': '*',
        'Access-Control-Allow-Methods': 'GET, POST, PATCH, PUT, DELETE, OPTIONS',
        'Access-Control-Allow-Headers': 'Origin, Content-Type, X-Auth-Token',
        'crossorigin': true,
        'Access-Control-Allow-Credentials': true,
      }
    ).then((response) => {

      // const url2 = window.URL.createObjectURL(new Blob([blob]),);
      console.log(response.data);
    });

  }

  // function populate_own_table() {
  //   var table = document.getElementById('own');

  //   var tr = document.createElement('tr');
  //   tr.innerHTML = '<td>hello</td>';
  //   table.appendChild(tr);


  // }


  // ShareFile();
  //deleteFile();
  // populate_own_table();
  return (
    <div className="flex flex-col justify-center align-center p-2 lg:pt-10 mx-auto max-w-[1440px] w-screen h-screen">
      {/* <h1 className="text-center font-bold text-xl">Welscome User2</h1> */}
      <div className="username_list_dropdown">
        <button onClick={handleDropdownOpen}>{username}</button>
        {openDropdown ? (
          <ul className="username_list_dropdown_menu">
            <li className="username_list_dropdown_menu-item">
              <button onClick={handleUsernameChange}>usr1</button>
            </li>
            <li className="username_list_dropdown_menu-item">
              <button onClick={handleUsernameChange}>usr2</button>
            </li>

          </ul>
        ) : null}
      </div>
      <div className="p-4 lg:pt-10 w-full h-full flex flex-col lg:flex-row justify-center align-center">
        <div className="lg:w-1/3">
          <form
            action=""
            method="post"
            encType="multipart/form-data"
            onSubmit={handleSubmit}
          >

            <h4 className="font-bold py-4">Upload File</h4>
            <label className="py-2 mt-2" htmlFor="button1">
              File to upload:
            </label>
            <br />
            <input type="file" onChange={handleChange} className="" />

            <br />
            <label className="py-2 mt-6" htmlFor="replication_factor">
              Replication Factor(between 1 to 100) &nbsp;
            </label>
            <input
              type="number"
              id="replication"
              className="border-[1px] border-[#000000]  max-w-[50px]"
            ></input>
            <br />
            <button
              onClick={async (e) => { }}
              type="submit"
              className="w-full md:max-w-[100px] mt-2 h-[28px] px-[10px] py-[6px] rounded-[4px] text-[8px] lg:text-[12px] font-bold bg-[#B1EAC1] hover:bg-[#58d17a]"
            >
              Upload
            </button>
          </form >
          <h4 className="font-bold py-4 lg:pt-10 ">Download File</h4>
          <button
            onClick={downloadFile}
            className="w-full md:max-w-[100px] mt-2 h-[28px] px-[10px] py-[6px] rounded-[4px] text-[8px] lg:text-[12px] font-bold bg-[#B1EAC1] hover:bg-[#58d17a]"
          >
            Download
          </button>
          <h4 className="font-bold py-4 lg:pt-10 ">Share File</h4>
          <form
            action=""
            method="post"
            encType="multipart/form-data"
            onSubmit={ShareFile}
          >
            <label htmlFor="share_with">Share With</label>
            <br />
            <input
              type="text"
              id="shareNameId"
              className="border-[1px] border-[#000000] max-w-[200px]"
            />
            <br />
            <button
              onClick={async (e) => { }}
              type="submit"
              className="w-full md:max-w-[100px] mt-2 h-[28px] px-[10px] py-[6px] rounded-[4px] text-[8px] lg:text-[12px] font-bold bg-[#B1EAC1] hover:bg-[#58d17a]"
            >
              Share
            </button>
          </form>
          <h4 className="font-bold py-4 lg:pt-10 ">Delete File</h4>
          <button
            onClick={deleteFile}
            className="w-full md:max-w-[100px] mt-2 h-[28px] px-[10px] py-[6px] rounded-[4px] text-[8px] lg:text-[12px] font-bold bg-[#B1EAC1] hover:bg-[#58d17a]"
          >
            Delete
          </button>
        </div>
        <div className="lg:w-1/3">
          <h4 className="font-bold py-4">Files:</h4>
          <table id="own" className="w-100 max-w-[400px] py-4 table-auto border-separate border-spacing-2">
            <thead>
              <tr className="font-bold text-base">
                <th className="w-[100px] text-left"> FileName</th>
              </tr>
            </thead>
            <tbody>
              {owner_files.map((file, index) => <tr key={`ownerfile_${index}`}>
                <td onClick={(e) => {
                  // console.log({
                  //   'ownerName': 'user1',
                  //   'fileName': owner_files[index]
                  // })
                  setselected(prev => {
                    return {
                      'ownerName': username,
                      'fileName': owner_files[index]
                    }
                  })

                }}><button className="file_name_item">{file}</button></td>
              </tr>)}
            </tbody>
          </table>
        </div>
        <div className="lg:w-1/3">
          <h4 className="font-bold py-4">Shared with me:</h4>
          <table className="py-4 table-auto">
            <thead>

              <tr className="font-bold text-base">
                <th className="w-[100px] text-left">FileName</th>
                <th className="w-[100px] text-left">Shared By</th>
              </tr>

            </thead>
            <tbody>
              {shared_files.map((sharedfile, index) =>
                <tr key={`sharedfile_${index}`} onClick={(e) => {
                  console.log({
                    'ownerName': shared_files[index][1],
                    'fileName': shared_files[index][0]
                  });
                  setselected(prev => {
                    return {
                      'ownerName': shared_files[index][1],
                      'fileName': shared_files[index][0]
                    }
                  })
                }}>
                  <td><button className="file_name_item">{sharedfile[0]}</button></td>
                  <td>{sharedfile[1]}</td>
                </tr>)}

            </tbody>
          </table>
        </div>
      </div >
    </div >
  );
}

export default App;
