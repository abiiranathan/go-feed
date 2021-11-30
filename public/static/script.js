"use strict";

class FeedController {
  constructor() {
    this.nextPageButton = document.querySelector("#next");
    this.prevPageButton = document.querySelector("#prev");
    this.currentPage = document.querySelector("#current_page");
    this.searchForm = document.querySelector("#search_form");

    this.page = 1;
    this.pagination = undefined;
    this.direction = "";
    this.feed = document.querySelector(".feed");
    this.setUpListeners();

    const feedForm = new NewFeedForm(this);
    feedForm.setUpListeners();
  }

  setUpListeners() {
    this.prevPageButton.addEventListener("click", async () => {
      if (this.page <= 1) return;
      this.page = this.page - 1;
      this.direction = "";
      await this.getFeed();
    });

    this.nextPageButton.addEventListener("click", async () => {
      if (this.page >= this.pagination.pages) return;
      this.page = this.page + 1;
      this.direction = "right";
      await this.getFeed();
    });

    this.searchForm.addEventListener("submit", e => {
      e.preventDefault();

      const query = document.getElementById("queryInput").value;
      this.getFeed(query || "");
    });
  }

  async getFeed(query) {
    this.feed.innerHTML = "";

    let url = `/feeds?limit=2`;

    if (query) {
      url += `&query=${query}&page=1`;
    } else {
      url += `&page=${this.page}`;
    }

    const response = await fetch(url, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    const { data, pagination } = await response.json();
    this.pagination = pagination;

    this.showHtml(data);
    this.handlePagination();
  }

  showHtml(data) {
    data.forEach(feed => {
      this.createFeedCard(feed);
    });
  }

  createFeedCard(feed) {
    const feedCard = document.createElement("div");
    feedCard.className = `feed-card ${this.direction}`;

    feedCard.innerHTML = `
      <a href="${feed.link}" target="_blank" class="feed-card-link">
          <div class="feed-card-title">${feed.title}</div>
          ${feed.image ? `<img class="feed-image" src="${feed.image}" alt="${feed.title}">` : ""}
          <div class="feed-card-description">${feed.description}</div>
          <a class="feed-card-link" href="${feed.link}">Read more</a>
      </a>  
    `;

    this.feed.appendChild(feedCard);
  }

  handlePagination() {
    this.page = this.pagination.page;

    if (this.page === 1) {
      this.prevPageButton.style.display = "none";
    } else {
      this.prevPageButton.style.display = "block";
    }

    if (this.pagination.pages === this.page) {
      this.nextPageButton.style.display = "none";
    } else {
      this.nextPageButton.style.display = "block";
    }

    this.currentPage.innerHTML = `Page ${this.page} of ${this.pagination.pages}`;
  }
}

class NewFeedForm {
  constructor(controller) {
    this.form = document.getElementById("newfeed");
    this.title = document.getElementById("title");
    this.description = document.getElementById("description");
    this.image = document.getElementById("image");
    this.link = document.getElementById("link");

    this.controller = controller;
  }

  getFormData() {
    return {
      title: this.title.value,
      description: this.description.value,
      image: this.image.value,
      link: this.link.value,
    };
  }

  setUpListeners() {
    this.form.addEventListener("submit", e => {
      e.preventDefault();

      const data = this.getFormData();
      this.submit(data);
    });
  }

  async submit(data) {
    try {
      const response = await fetch("/feeds", {
        method: "POST",
        headers: {
          "content-type": "application/json",
          accepts: "application/json",
        },
        body: JSON.stringify(data),
      });

      await response.json();
      this.controller.getFeed();
      this.form.reset();
    } catch (error) {
      console.log(error);
      alert("Something went wrong");
    }
  }
}

const feedController = new FeedController();
feedController.getFeed();
