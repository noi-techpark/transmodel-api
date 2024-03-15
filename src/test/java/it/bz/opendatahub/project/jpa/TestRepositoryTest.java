// SPDX-FileCopyrightText: NOI Techpark <digital@noi.bz.it>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package it.bz.opendatahub.project.jpa;



import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;

import static org.junit.jupiter.api.Assertions.*;

@DataJpaTest
public class TestRepositoryTest {

    @Autowired
    private TestRepository testRepository;

    @Test
    public void givenEmptyDB_WhenInsertNewRow_ThanSave() {

        it.bz.opendatahub.project.model.Test model = new it.bz.opendatahub.project.model.Test();
        model.setText("Hello World");

        final var result = testRepository.save(model);

        assertEquals("Hello World", result.getText());
        assertEquals(1L, result.getId());
    }

}