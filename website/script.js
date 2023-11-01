//Checkport
/* <div class="wrapper">
            <div class="container">
              <form id="pingForm" action="/ping" method="POST">
                <div class="input-bar">
                  <div class="input-container">
                    <div class="input-wrapper">
                      <input type="text" id="ipInput" name="ip" class="input" placeholder="Enter IP address">
                    </div>
                    <div class="input-wrapper">
                      <input type="text" id="portInput" name="port" class="input" placeholder="Enter port">
                    </div>
                  </div>
                  <button type="button" class="input-btn" onclick="sendIPAndPort()">Port Check</button>
                </div>
              </form>
            </div>
          </div> */
function sendIPAndPort() {
    const ipInput = $("#ipInput");
    const portInput = $("#portInput");
    const textArea = $("#outputTextarea");
    
    if (ipInput.val() === "" || ipInput.val().trim() === "" || portInput.val() === "" || portInput.val().trim() === "") {
      textArea.val("Please enter both IP address and port");
      return;
    }
    textArea.val("");
    
    const formData = new FormData();
    formData.append("ip", ipInput.val());
    formData.append("port", portInput.val());
    
    $.ajax({
      url: "/checkport",
      type: "POST",
      data: formData,
      processData: false,
      contentType: false,
      success: function(response) {
        textArea.val(textArea.val() + response);
      },
      error: function(xhr, status, error) {
        textArea.val(textArea.val() + "Error: " + xhr.responseText);
      }
    });
    
    ipInput.val("");
    portInput.val("");
    ipInput.focus();
  }

  //ipinfo
  /* <div class="wrapper">
        <div class="container">
          <div class="input-bar">
            <input
              type="text"
              id="ipInput"
              name="ip"
              class="input"
              placeholder="Enter IP address"
            />
            <button type="button" class="input-btn" onclick="sendIP()">info</button>
          </div>
        </div>
      </div>*/
  function sendIP() {
    const ipInput = $("#ipInput");
    const textArea = $("#outputTextarea");
  
    if (ipInput.val() === "" || ipInput.val().trim() === "") {
      textArea.val("Please enter an IP address");
      return;
    }
    textArea.val("");
  
    const formData = new FormData();
    formData.append("ip", ipInput.val());
  
    $.ajax({
      url: "/location",
      type: "POST",
      data: formData,
      processData: false,
      contentType: false,
      dataType: "json", 
      success: function(response) {
        textArea.val(JSON.stringify(response, null, 4));
      },
      error: function(xhr, status, error) {
        textArea.val("Error: " + xhr.responseText);
      }
    });
  
    ipInput.val("");
    ipInput.focus();
  }
  
  
  //Ping
  /*
   <div class="wrapper">
          <div class="container">
            <div class="input-bar">
              <input type="text" id="ipInput" name="ip" class="input" placeholder="Enter IP address">
              <button type="button" class="input-btn" onclick="sendIP()">Ping</button>
            </div>
          </div>
        </div>
   */
  function sendIP() {
    const ipInput = document.getElementById("ipInput");
    const textArea = document.getElementById("outputTextarea");
    if (ipInput.value === "" || ipInput.value.trim() === "") {
      textArea.value = "Please enter an IP address";
      return;
    }
    textArea.value = "";

    const formData = new FormData();
    formData.append("ip", ipInput.value);

    const xhr = new XMLHttpRequest();
    xhr.open("POST", "/ping", true);
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xhr.onreadystatechange = function () {
      if (xhr.readyState === XMLHttpRequest.DONE) {
        if (xhr.status === 200) {
          const response = JSON.parse(xhr.responseText);
          textArea.value += response.result;
        } else {
          textArea.value += "Error: " + xhr.responseText;
        }
      }
    };
    xhr.send(new URLSearchParams(formData));

    ipInput.value = "";
    ipInput.focus();
  }

  //traceroute

  /*
   <div class="flex-container">
        <div>
          <div class="wrapper">
            <div class="container">
              <div class="input-bar">
                <input type="text" id="ipInput" name="ip" class="input"
                  placeholder="Enter IP address" />
                <button type="button" class="input-btn" onclick="sendIP()">Traceroute</button>
              </div>
            </div>
          </div>
          <div class="output">
            <textarea id="outputTextarea"></textarea>
          </div>
        </div>
      </div>*/
  function sendIP() {
    const ipInput = $("#ipInput");
    const textArea = $("#outputTextarea");
  
    if (ipInput.val() === "" || ipInput.val().trim() === "") {
      textArea.val("Please enter an IP address");
      return;
    }
    textArea.val("");
  
    const formData = new FormData();
    formData.append("ip", ipInput.val());
  
    $.ajax({
      url: "/traceroute",
      type: "POST",
      data: new URLSearchParams(formData).toString(),
      contentType: "application/x-www-form-urlencoded",
      success: function(response) {
        textArea.val(textArea.val() + response);
      },
      error: function(xhr, status, error) {
        textArea.val(textArea.val() + "Error: " + xhr.responseText);
      }
    });
  
    ipInput.val("");
    ipInput.focus();
  }
  
  //whois
  /*<div class="flex-container">
        <div>
          <div class="wrapper">
            <div class="container">
              <form id="pingForm" action="/whois" method="POST">
                <div class="input-bar">
                  <input type="text" id="ipInput" name="ip" class="input"
                    placeholder="Enter IP address">
                  <button type="button" class="input-btn" onclick="sendIP()">Who
                    is</button>
                </div>
              </form>

            </div>
          </div>
          <div class="output">
            <textarea id="outputTextarea"></textarea>
          </div>
        </div>
      </div>*/
  function sendIP() {
    const ipInput = $("#ipInput");
    const textArea = $("#outputTextarea");
  
    textArea.val("");
  
    const formData = new FormData();
    formData.append("ip", ipInput.val());
  
    $.ajax({
      url: "/WHOIS",
      type: "POST",
      data: formData,
      processData: false,
      contentType: false,
      success: function(response) {
        textArea.val(response.result);
      },
      error: function(xhr, status, error) {
        textArea.val("Error: " + xhr.responseText);
      }
    });
  
    ipInput.val("");
    ipInput.focus();
  }
  