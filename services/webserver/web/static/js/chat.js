var conn = undefined;

function disableChat() {
  $(document).ready(function () {
    $("#message-to-send").off().on("keydown", function (event) {
      event.preventDefault();
      return;
    });
  });
}

function renderError(errorMessage) {
  var messageElement = document.createElement('div');
  messageElement.className = "message";
  messageElement.innerText = errorMessage;
  $('.chat-history').find('ul').append(messageElement);
  toastr.error(errorMessage, 'Chat');
}

function renderGlobal(message){
  toastr.info(message, 'Global Notification', {timeOut: 10000});
}

function chatInterface() {
  var chat = {
    init: function () {
      this.cacheDOM();
      this.bindEvents();
    },
    cacheDOM: function () {
      if ($("#chat").length){
        this.chatHistory = $('.chat-history');
        this.button = $('#send-button');
        this.textarea = $('#message-to-send');
        this.chatHistoryList = this.chatHistory.find('ul');
      }
    },
    bindEvents: function () {
      if ($("#chat").length){
        this.button.on('click', this.sendMessage.bind(this));
        this.textarea.on('keyup', this.sendMessageEnter.bind(this));
      }
    },
    render: function (message) {
      if ($("#chat").length && message.Data != '') {
        var messageElement = document.createElement('div');
        messageElement.className = "message";

        if (message.Auto){
            messageElement.innerText = `[${message.Time}] ${message.Data}`;
        } else {
            messageElement.innerText = `[${message.Time}] ${message.Name}: ${message.Data}`;
        }

        this.chatHistoryList.append(messageElement);
        this.scrollToBottom();
      }
    },
    sendMessage: function () {
      var message = this.textarea.val()
      this.textarea.val('');
      if (typeof conn != "undefined" && message.trim() !== '') {
        conn.send(JSON.stringify({
          "Data": message
        }));
      }
    },
    sendMessageEnter: function (event) {
      if (event.keyCode === 13) {
        this.sendMessage();
      }
    },
    scrollToBottom: function () {
      this.chatHistory.scrollTop(this.chatHistory[0].scrollHeight);
    },
  };

  if (window["WebSocket"]) {
    conn = new WebSocket("ws://" + document.location.host + "/chat/ws");

    conn.onclose = function (evt) {
      if (evt.code == 1006) {
        renderError("Could not connect to the chat server.");
      } else {
        renderError("Chat server error. Please reload the page.");
      }

      conn = null;
      disableChat();
    };

    conn.onmessage = function (event) {
      const json = (function (resp) {
        try {
          return JSON.parse(resp);
        } catch (err) {
          return false;
        }
      })(event.data);

      if (!json) {
        toastr.error('Action Request Failed.', 'Chat');
        return;
      }

      if (json.Global){
        renderGlobal(json.Data);
      }

      if (json.Messages) {
        json.Messages.forEach(el => {
          chat.render(el);
        });
      } else {
        chat.render(json);
      }
    };

    $(window).on('beforeunload', function () {
      if (conn) {
        conn.close();
      }
    });

    chat.init();
  } else {
    disableChat();
    renderError("Your browser does not support Websockets.");
  }
}

chatInterface();
