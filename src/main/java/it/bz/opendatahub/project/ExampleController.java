// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package it.bz.opendatahub.project;

import it.bz.opendatahub.project.jpa.TestRepository;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import static net.logstash.logback.argument.StructuredArguments.v;

@RestController
public class ExampleController {

    Logger logger = LoggerFactory.getLogger(ExampleController.class);

    private final TestRepository testRepository;

    public ExampleController(TestRepository testRepository) {
        this.testRepository = testRepository;
    }

    @RequestMapping("/")
    public String index() {
        logger.info("We can log an ID as key-value pair for easy lookup in Log storage such as Elasticsearch.", v("id", 13));
        testRepository.findAll();
        return "Hello World";
    }
}
